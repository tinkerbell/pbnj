package oob

import (
	"context"

	"github.com/bmc-toolbox/bmclib/devices"
	"github.com/bmc-toolbox/bmclib/discover"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
)

// User management methods
type User interface {
	Create(ctx context.Context, in v1.CreateUserRequest) (result string, err *Error)
	Update(ctx context.Context, in v1.UpdateUserRequest) (result string, err *Error)
	Delete(ctx context.Context, in v1.DeleteUserRequest) (result string, err *Error)
}

// Machine management methods
type Machine interface {
	BootDevice(ctx context.Context, in *v1.DeviceRequest) (result string, err *Error)
	Power(ctx context.Context, status chan string, in *v1.PowerRequest) (result string, err *Error)
}

// BMC management methods
type BMC interface {
	Reset(ctx context.Context, in v1.ResetRequest) (result string, err *Error)
	NetworkSource(ctx context.Context, in v1.NetworkSourceRequest) (result string, err *Error)
}

// Conn details for connecting
type Conn struct {
	Log logging.Logger
}

// Error for all bmc actions
type Error struct {
	*v1.Error
}

// BootDevice functionality for machines
func (c *Conn) BootDevice(ctx context.Context, in *v1.DeviceRequest) (result string, err *Error) {
	l := c.Log.GetContextLogger(ctx)
	l.V(0).Info("not implemented")
	return result, err
}

// Power functionality for machines
func (c *Conn) Power(ctx context.Context, status chan string, in *v1.PowerRequest) (string, *Error) {
	l := c.Log.GetContextLogger(ctx)
	l.V(0).Info("power state")
	// TODO handle nil values
	var result string
	errMsg := &Error{
		Error: &v1.Error{
			Code:    0,
			Message: "",
			Details: nil,
		},
	}

	if in.Authn == nil || in.Authn.Authn == nil {
		msg := "no auth found"
		errMsg.Code = 16
		errMsg.Message = msg
		return msg, errMsg
	}
	user := in.GetAuthn().GetDirectAuthn().Username
	password := in.GetAuthn().GetDirectAuthn().Password
	host := in.GetAuthn().GetDirectAuthn().GetHost().Host

	status <- "trying to connect to bmc"

	connection, err := discover.ScanAndConnect(host, user, password, discover.WithLogger(l))
	if err != nil {
		// TODO set errMsg.Code based on err response
		errMsg.Code = 2
		errMsg.Message = err.Error()
		return result, errMsg
	}
	status <- "connected to bmc"

	switch connection := connection.(type) {

	case devices.Bmc:
		conn := connection.(devices.Bmc)
		defer conn.Close()

		switch in.GetAction().String() {
		case v1.PowerRequest_ON.String():
			// ok, err := conn.PowerOn()
			msg := "power ON not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_OFF.String():
			// ok, err := conn.PowerOff()
			msg := "power OFF not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_HARDOFF.String():
			msg := "power HARD OFF not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_CYCLE.String():
			// ok, err := conn.PowerCycle()
			msg := "power CYCLE not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_RESET.String():
			msg := "power RESET not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_STATUS.String():
			l.V(1).Info("getting power status")
			status <- "getting power status"
			result, err := conn.PowerState()
			if err != nil {
				// TODO need to set code based on response
				status <- "error getting power state"
				errMsg.Code = 2
				errMsg.Message = err.Error()
			}
			return result, errMsg
		}

	case devices.Cmc:
		l.V(1).Info("type cmc detected")
		l.V(0).Info("not implemented")
		conn := connection.(devices.Cmc)

		switch in.GetAction().String() {
		case v1.PowerRequest_ON.String():
			// ok, err := conn.PowerOn()
			msg := "power ON not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_OFF.String():
			// ok, err := conn.PowerOff()
			msg := "power OFF not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_HARDOFF.String():
			msg := "power HARD OFF not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_CYCLE.String():
			// ok, err := conn.PowerCycle()
			msg := "power CYCLE not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_RESET.String():
			msg := "power RESET not implemented"
			l.V(1).Info(msg)
			errMsg.Code = 12
			errMsg.Message = msg
			return result, errMsg
		case v1.PowerRequest_STATUS.String():
			l.V(0).Info("getting power status")
			status <- "getting power status"
			result, err := conn.Status()
			if err != nil {
				// TODO need to set code based on response
				status <- "error getting power state"
				errMsg.Code = 2
				errMsg.Message = err.Error()
			}
			return result, errMsg
		}

	default:
		errMsg.Code = 12
		errMsg.Message = "Unknown device"
		return result, errMsg
	}

	return result, errMsg

}
