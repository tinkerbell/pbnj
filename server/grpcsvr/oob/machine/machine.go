package machine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
	bmc "github.com/tinkerbell/pbnj/server/grpcsvr/oob"
)

// Action for making power actions on BMCs, implements oob.Machine interface
type Action struct {
	bmc.Accessory
	PowerRequest      *v1.PowerRequest
	BootDeviceRequest *v1.DeviceRequest
}

type powerConnection struct {
	bmc.ConnectionDetails
	pwr power
}

type bootDeviceConnection struct {
	bmc.ConnectionDetails
	boot
}

// the power interface allows us to abstract these functions
// between different libraries and BMC connections
// like ipmi, racadm, redfish, etc
type power interface {
	bmc.Connection
	on() (string, repository.Error)
	off() (string, repository.Error)
	status() (string, repository.Error)
	reset() (string, repository.Error)
	hardoff() (string, repository.Error)
	cycle() (string, repository.Error)
}

// the boot interface allows us to abstract these functions
// between different libraries and BMC connections
// like ipmi, racadm, redfish, etc
type boot interface {
	bmc.Connection
	setBootDevice() (string, repository.Error)
}

// Option to add to an Actions
type Option func(a *Action) error

// WithContext adds a context to an Action struct
func WithContext(c context.Context) Option {
	return func(a *Action) error {
		a.Ctx = c
		return nil
	}
}

// WithLogger adds a logr to an Action struct
func WithLogger(l logr.Logger) Option {
	return func(a *Action) error {
		a.Log = l
		return nil
	}
}

// WithStatusMessage adds a status message chan to an Action struct
func WithStatusMessage(s chan string) Option {
	return func(a *Action) error {
		a.StatusMessages = s
		return nil
	}
}

// WithDeviceRequest adds DeviceRequest to an Action struct
func WithDeviceRequest(in *v1.DeviceRequest) Option {
	return func(a *Action) error {
		a.BootDeviceRequest = in
		return nil
	}
}

// WithPowerRequest adds PowerRequest to an Action struct
func WithPowerRequest(in *v1.PowerRequest) Option {
	return func(a *Action) error {
		a.PowerRequest = in
		return nil
	}
}

// NewMachine returns an oob.Machine interface
func NewMachine(opts ...Option) (oob.Machine, error) {
	a := &Action{}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// BootDevice functionality for machines
func (m Action) BootDevice(ctx context.Context) (result string, errMsg repository.Error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	host, user, password, errMsg := m.ParseAuth(m.BootDeviceRequest.Authn)
	if errMsg.Message != "" {
		return result, errMsg
	}

	base := "setting boot device: " + m.BootDeviceRequest.GetDevice().String()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := []bootDeviceConnection{
		{ConnectionDetails: bmc.ConnectionDetails{Name: "ipmi"}, boot: &ipmiBootDevice{mAction: m, user: user, password: password, host: host, port: "623"}},
	}

	var connected bool
	m.SendStatusMessage("connecting to BMC")
	for index := range connections {
		connections[index].Err = connections[index].Connect(ctx)
		if connections[index].Err.Message == "" {
			connections[index].Connected = true
			defer connections[index].Close()
			connected = true
		}
	}
	m.Log.V(1).Info("connections", "connections", fmt.Sprintf("%+v", connections))
	if !connected {
		m.SendStatusMessage("connecting to BMC failed")
		var combinedErrs []string
		for _, connection := range connections {
			combinedErrs = append(combinedErrs, connection.Err.Message)
		}
		msg := "could not connect"
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = msg
		errMsg.Details = append(errMsg.Details, combinedErrs...)
		m.Log.V(0).Info(msg, "error", combinedErrs)
		return result, errMsg
	}
	m.SendStatusMessage("connected to BMC")

	for index := range connections {
		if connections[index].Connected {
			m.Log.V(0).Info("trying", "name", connections[index].Name)
			result, errMsg = connections[index].setBootDevice()
			if errMsg.Message == "" {
				m.Log.V(0).Info("action implemented by", "implementer", connections[index].Name)
				break
			}
		}
	}

	if errMsg.Message != "" {
		m.SendStatusMessage("error with " + base + ": " + errMsg.Message)
		m.Log.V(0).Info("error with "+base, "error", errMsg.Message)
	}
	m.SendStatusMessage(base + " complete")
	return strings.ToLower(result), errMsg //nolint
}

// Power functionality for machines
func (m Action) Power(ctx context.Context) (result string, errMsg repository.Error) {
	host, user, password, errMsg := m.ParseAuth(m.PowerRequest.Authn)
	if errMsg.Message != "" {
		return result, errMsg
	}

	base := "power " + m.PowerRequest.GetAction().String()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	// the order here is the order in which these connections/operations will be tried
	connections := []powerConnection{
		{ConnectionDetails: bmc.ConnectionDetails{Name: "bmclib"}, pwr: &bmclibBMC{user: user, password: password, host: host, log: m.Log}},
		{ConnectionDetails: bmc.ConnectionDetails{Name: "ipmi"}, pwr: &ipmiBMC{user: user, password: password, host: host, ctx: m.Ctx, log: m.Log}},
		{ConnectionDetails: bmc.ConnectionDetails{Name: "redfish"}, pwr: &redfishBMC{user: user, password: password, host: host, log: m.Log}},
	}

	var connected bool
	m.SendStatusMessage("connecting to BMC")
	for index := range connections {
		connections[index].Err = connections[index].pwr.Connect(ctx)
		if connections[index].Err.Message == "" {
			connections[index].Connected = true
			defer connections[index].pwr.Close()
			connected = true
		}
	}
	m.Log.V(1).Info("connections", "connections", fmt.Sprintf("%+v", connections))
	if !connected {
		m.SendStatusMessage("connecting to BMC failed")
		var combinedErrs []string
		for _, connection := range connections {
			combinedErrs = append(combinedErrs, connection.Err.Message)
		}
		msg := "could not connect"
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = msg
		errMsg.Details = append(errMsg.Details, combinedErrs...)
		m.Log.V(0).Info(msg, "error", combinedErrs)
		return result, errMsg
	}
	m.SendStatusMessage("connected to BMC")

	for _, connection := range connections {
		if connection.Connected {
			m.Log.V(1).Info("trying", "name", connection.Name)
			result, errMsg = doAction(m.PowerRequest.GetAction(), connection.pwr)
			if errMsg.Message == "" {
				m.Log.V(1).Info("action implemented by", "implementer", connection.Name)
				break
			}
		}
	}

	if errMsg.Message != "" {
		m.SendStatusMessage("error with " + base + ": " + errMsg.Message)
		m.Log.V(0).Info("error with "+base, "error", errMsg.Message)
	}
	m.SendStatusMessage(base + " complete")
	return strings.ToLower(result), errMsg //nolint
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
