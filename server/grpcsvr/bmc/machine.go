package bmc

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// MachineAction for making power actions on BMCs, implements oob.Machine interface
type MachineAction struct {
	Log               logging.Logger
	Ctx               context.Context
	PowerRequest      *v1.PowerRequest
	BootDeviceRequest *v1.DeviceRequest
	StatusMessages    chan string
}

type power interface {
	connection() repository.Error
	close()
	on() (string, repository.Error)
	off() (string, repository.Error)
	status() (string, repository.Error)
	reset() (string, repository.Error)
	hardoff() (string, repository.Error)
	cycle() (string, repository.Error)
}

// BootDevice functionality for machines
func (m MachineAction) BootDevice() (result string, errMsg repository.Error) {
	l := m.Log.GetContextLogger(m.Ctx)
	msg := "setting Boot Device not implemented yet"
	l.V(0).Info(msg)
	errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
	errMsg.Message = msg
	return result, errMsg //nolint
}

type powerConnection struct {
	name      string
	connected bool
	pwr       power
	err       repository.Error
}

// Power functionality for machines
func (m MachineAction) Power() (result string, errMsg repository.Error) {
	l := m.Log.GetContextLogger(m.Ctx)
	host, user, password, errMsg := m.parseAuth(m.PowerRequest.Authn)
	if errMsg.Message != "" {
		return result, errMsg
	}

	base := "power " + m.PowerRequest.GetAction().String()
	msg := "working on " + base
	m.sendStatusMessage(msg)

	// the order here is the order in which these connections/operations will be tried
	connections := []powerConnection{
		{name: "bmclib", pwr: &bmclibBMC{mAction: m, user: user, password: password, host: host}},
		{name: "ipmi", pwr: &ipmiBMC{mAction: m, user: user, password: password, host: host}},
		{name: "redfish", pwr: &redfishBMC{mAction: m, user: user, password: password, host: host}},
	}

	var connected bool
	m.sendStatusMessage("connecting to BMC")
	for index := range connections {
		connections[index].err = connections[index].pwr.connection()
		if connections[index].err.Message == "" {
			connections[index].connected = true
			defer connections[index].pwr.close()
			connected = true
		}
	}
	l.V(1).Info("connections", "connections", fmt.Sprintf("%+v", connections))
	if !connected {
		m.sendStatusMessage("connecting to BMC failed")
		var combinedErrs []string
		for _, connection := range connections {
			combinedErrs = append(combinedErrs, connection.err.Message)
		}
		msg := "could not connect"
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = msg
		errMsg.Details = append(errMsg.Details, combinedErrs...)
		l.V(0).Info(msg, "error", combinedErrs)
		return result, errMsg
	}
	m.sendStatusMessage("connected to BMC")

	for _, connection := range connections {
		if connection.connected {
			l.V(1).Info("trying", "name", connection.name)
			result, errMsg = doAction(m.PowerRequest.GetAction(), connection.pwr)
			if errMsg.Message == "" {
				l.V(1).Info("action implemented by", "implementer", connection.name)
				break
			}
		}
	}

	if errMsg.Message != "" {
		m.sendStatusMessage("error with " + base + ": " + errMsg.Message)
		l.V(0).Info("error with "+base, "error", errMsg.Message)
	}
	m.sendStatusMessage(base + " complete")
	return result, errMsg //nolint
}

func doAction(action v1.PowerRequest_Action, pwr power) (result string, errMsg repository.Error) {
	switch action {
	case v1.PowerRequest_ON:
		result, errMsg = pwr.on()
	case v1.PowerRequest_OFF:
		result, errMsg = pwr.off()
	case v1.PowerRequest_STATUS:
		result, errMsg = pwr.status()
	case v1.PowerRequest_RESET:
		result, errMsg = pwr.reset()
	case v1.PowerRequest_HARDOFF:
		result, errMsg = pwr.hardoff()
	case v1.PowerRequest_CYCLE:
		result, errMsg = pwr.cycle()
	default:
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "unknown power action"
	}
	return result, errMsg
}

func (m MachineAction) sendStatusMessage(msg string) {
	select {
	case m.StatusMessages <- msg:
		return
	case <-time.After(2 * time.Second):
		l := m.Log.GetContextLogger(m.Ctx)
		l.V(0).Info("timed out waiting for status message receiver", "statusMsg", msg)
	}
}

func (m MachineAction) parseAuth(auth *v1.Authn) (host string, username string, passwd string, errMsg repository.Error) {
	if auth == nil || auth.Authn == nil || auth.GetDirectAuthn() == nil {
		msg := "no auth found"
		m.sendStatusMessage(msg)
		errMsg.Code = v1.Code_value["UNAUTHENTICATED"]
		errMsg.Message = msg
		return
	}

	username = auth.GetDirectAuthn().GetUsername()
	passwd = auth.GetDirectAuthn().GetPassword()
	host = auth.GetDirectAuthn().GetHost().GetHost()

	return host, username, passwd, errMsg
}
