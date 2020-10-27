package bmc

import (
	"context"

	"github.com/bmc-toolbox/bmclib/devices"
	"github.com/bmc-toolbox/bmclib/discover"
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

// BootDevice functionality for machines
func (p MachineAction) BootDevice() (string, repository.Error) {
	var result string
	errMsg := repository.Error{
		Code:    0,
		Message: "",
		Details: nil,
	}
	l := p.Log.GetContextLogger(p.Ctx)
	l.V(0).Info("not implemented")
	msg := "power OFF not implemented"
	l.V(1).Info(msg)
	errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
	errMsg.Message = msg
	return result, errMsg //nolint
}

// Power functionality for machines
func (p MachineAction) Power() (string, repository.Error) {
	l := p.Log.GetContextLogger(p.Ctx)
	l.V(0).Info("power state")
	// TODO handle nil values
	var result string
	errMsg := repository.Error{
		Code:    0,
		Message: "",
		Details: nil,
	}

	if p.PowerRequest.Authn == nil || p.PowerRequest.Authn.Authn == nil {
		msg := "no auth found"
		errMsg.Code = v1.Code_value["UNAUTHENTICATED"]
		errMsg.Message = msg
		return msg, errMsg //nolint
	}
	user := p.PowerRequest.GetAuthn().GetDirectAuthn().Username
	password := p.PowerRequest.GetAuthn().GetDirectAuthn().Password
	host := p.PowerRequest.GetAuthn().GetDirectAuthn().GetHost().Host

	p.StatusMessages <- "trying to connect to bmc"

	connection, err := discover.ScanAndConnect(host, user, password, discover.WithLogger(l))
	if err != nil {
		// TODO set errMsg.Code based on err response
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return result, errMsg //nolint
	}
	p.StatusMessages <- "connected to bmc"

	switch connection := connection.(type) {

	case devices.Bmc:
		conn := connection.(devices.Bmc)
		defer conn.Close()

		switch p.PowerRequest.GetAction() {
		case v1.PowerRequest_ON:
			// ok, err := conn.PowerOn()
			msg := "power ON not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_OFF:
			// ok, err := conn.PowerOff()
			msg := "power OFF not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_HARDOFF:
			msg := "power HARD OFF not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_CYCLE:
			// ok, err := conn.PowerCycle()
			msg := "power CYCLE not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_RESET:
			msg := "power RESET not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_STATUS:
			l.V(1).Info("getting power status")
			p.StatusMessages <- "getting power status"
			result, err := conn.PowerState()
			if err != nil {
				// TODO need to set code based on response
				p.StatusMessages <- "error getting power state"
				errMsg.Code = v1.Code_value["UNKNOWN"]
				errMsg.Message = err.Error()
			}
			return result, errMsg //nolint
		}

	case devices.Cmc:
		l.V(1).Info("type cmc detected")
		l.V(0).Info("not implemented")
		conn := connection.(devices.Cmc)

		switch p.PowerRequest.GetAction() {
		case v1.PowerRequest_ON:
			// ok, err := conn.PowerOn()
			msg := "power ON not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_OFF:
			// ok, err := conn.PowerOff()
			msg := "power OFF not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_HARDOFF:
			msg := "power HARD OFF not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_CYCLE:
			// ok, err := conn.PowerCycle()
			msg := "power CYCLE not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_RESET:
			msg := "power RESET not implemented"
			l.V(1).Info(msg)
			errMsg.Code = v1.Code_value["UNIMPLEMENTED"]
			errMsg.Message = msg
			return result, errMsg //nolint
		case v1.PowerRequest_STATUS:
			l.V(0).Info("getting power status")
			p.StatusMessages <- "getting power status"
			result, err := conn.Status()
			if err != nil {
				// TODO need to set code based on response
				p.StatusMessages <- "error getting power state"
				errMsg.Code = 2
				errMsg.Message = err.Error()
			}
			return result, errMsg //nolint
		}

	default:
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "Unknown device"
		return result, errMsg //nolint
	}

	return result, errMsg //nolint

}
