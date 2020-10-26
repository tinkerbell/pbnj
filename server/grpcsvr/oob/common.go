package oob

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// ConnectionDetails shared amongst all implementations
type ConnectionDetails struct {
	Name      string
	Connected bool
	Err       repository.Error
}

// Connection methods open/close
type Connection interface {
	Connect(context.Context) repository.Error
	Close()
}

// Accessory for all BMC actions
type Accessory struct {
	Log            logr.Logger
	Ctx            context.Context
	StatusMessages chan string
}

// ParseAuth will return host, user, passwd from auth struct
func (a *Accessory) ParseAuth(auth *v1.Authn) (host string, username string, passwd string, errMsg repository.Error) {
	if auth == nil || auth.Authn == nil || auth.GetDirectAuthn() == nil {
		msg := "no auth found"
		a.SendStatusMessage(msg)
		errMsg.Code = v1.Code_value["UNAUTHENTICATED"]
		errMsg.Message = msg
		return
	}

	username = auth.GetDirectAuthn().GetUsername()
	passwd = auth.GetDirectAuthn().GetPassword()
	host = auth.GetDirectAuthn().GetHost().GetHost()

	return host, username, passwd, errMsg
}

// SendStatusMessage will send a message to a string chan
func (a *Accessory) SendStatusMessage(msg string) {
	select {
	case a.StatusMessages <- msg:
		return
	case <-time.After(2 * time.Second):
		a.Log.V(1).Info("timed out waiting for status message receiver", "statusMsg", msg)
	}
}
