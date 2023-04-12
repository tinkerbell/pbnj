package oob

import (
	"context"
	"errors"
	"time"

	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

var (
	// DefaultBmcTimeout is the default value for how long a BMC call/interaction is allowed to run before it is cancelled.
	DefaultBmcTimeout = 120 * time.Second
)

// Connection methods open/close.
type Connection interface {
	Connect(context.Context) error
	Close(context.Context)
}

// Accessory for all BMC actions.
type Accessory struct {
	Log            logr.Logger
	StatusMessages chan string
	// SkipRedfishVersions is a list of Redfish versions to be ignored,
	//
	// When running an action on a BMC, PBnJ will pass the value of the skipRedfishVersions to bmclib
	// which will then ignore the Redfish endpoint completely on BMCs running the given Redfish versions,
	// and will proceed to attempt other drivers like - IPMI/SSH/Vendor API instead.
	//
	// for more information see https://github.com/bmc-toolbox/bmclib#bmc-connections
	SkipRedfishVersions []string
}

// Connect to a BMC interface function.
func Connect(ctx context.Context, conn Connection) error {
	return conn.Connect(ctx)
}

// Close a BMC interface function.
func Close(ctx context.Context, conn Connection) {
	conn.Close(ctx)
}

// EstablishConnections tries to connect to all BMCs.
// Successful connection names are returned in a slice of strings.
// If no connections were successful then an error is returned.
func EstablishConnections(ctx context.Context, bmcs map[string]interface{}) (successfulConnections []string, err error) {
	var connErrs []error
	var connected bool

	for name, elem := range bmcs {
		switch con := elem.(type) {
		case Connection:
			connErr := Connect(ctx, con)
			if connErr == nil {
				successfulConnections = append(successfulConnections, name)
				connected = true
			} else {
				connErrs = append(connErrs, connErr)
			}
		default:
			connErrs = append(connErrs, errors.New("unknown connection type"))
		}
	}

	if !connected {
		var combinedErrs []string
		for _, connection := range connErrs {
			if connection != nil {
				combinedErrs = append(combinedErrs, connection.Error())
			}
		}
		errMsg := repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "could not connect",
			Details: combinedErrs,
		}
		return successfulConnections, &errMsg
	}
	return successfulConnections, nil
}

// ParseAuth will return host, user, passwd from auth struct.
func (a *Accessory) ParseAuth(auth *v1.Authn) (host string, username string, passwd string, err error) {
	var errMsg repository.Error
	if auth == nil || auth.Authn == nil || auth.GetDirectAuthn() == nil {
		msg := "no auth found"
		a.SendStatusMessage(msg)
		errMsg.Code = v1.Code_value["UNAUTHENTICATED"]
		errMsg.Message = msg
		return host, username, passwd, &errMsg
	}

	username = auth.GetDirectAuthn().GetUsername()
	passwd = auth.GetDirectAuthn().GetPassword()
	host = auth.GetDirectAuthn().GetHost().GetHost()

	return host, username, passwd, nil
}

// SendStatusMessage will send a message to a string chan.
func (a *Accessory) SendStatusMessage(msg string) {
	select {
	case a.StatusMessages <- msg:
		return
	case <-time.After(2 * time.Second):
		a.Log.V(1).Info("timed out waiting for status message receiver", "statusMsg", msg)
	}
}
