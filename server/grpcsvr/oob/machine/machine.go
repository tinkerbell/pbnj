package machine

import (
	"context"

	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/oob"
	common "github.com/tinkerbell/pbnj/server/grpcsvr/oob"
)

// Action for making power actions on BMCs, implements oob.Machine interface
type Action struct {
	common.Accessory
	PowerRequest      *v1.PowerRequest
	BootDeviceRequest *v1.DeviceRequest
}

// Option to add to an Actions
type Option func(a *Action) error

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
func (m Action) BootDevice(ctx context.Context, device string) (result string, err error) {
	host, user, password, parseErr := m.ParseAuth(m.BootDeviceRequest.Authn)
	if parseErr != nil {
		return result, parseErr
	}
	base := "setting boot device: " + m.BootDeviceRequest.GetBootDevice().String()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := []interface{}{
		&ipmiBootDevice{mAction: m, user: user, password: password, host: host, port: "623"},
	}

	m.SendStatusMessage("connecting to BMC")
	successfulConnections, ecErr := common.EstablishConnections(ctx, connections)
	if ecErr != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return result, ecErr
	}
	m.SendStatusMessage("connected to BMC")

	var userAction []oob.Machine
	for _, elem := range successfulConnections {
		elem := *elem
		switch r := elem.(type) {
		case oob.Machine:
			userAction = append(userAction, r)
		}
	}
	result, err = oob.MachineBootDevice(ctx, device, userAction)
	if err != nil {
		m.SendStatusMessage("error with " + base + ": " + err.Error())
		m.Log.V(0).Info("error with "+base, "error", err.Error())
		return result, err
	}
	m.SendStatusMessage(base + " complete")
	return result, nil
}

// Power functionality for machines
func (m Action) Power(ctx context.Context, action string) (result string, err error) {

	host, user, password, parseErr := m.ParseAuth(m.PowerRequest.Authn)
	if parseErr != nil {
		return result, parseErr
	}
	base := "power " + m.PowerRequest.GetPowerAction().String()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	// the order here is the order in which these connections/operations will be tried
	connections := []interface{}{
		&bmclibBMC{user: user, password: password, host: host, log: m.Log},
		&ipmiBMC{user: user, password: password, host: host, log: m.Log},
		&redfishBMC{user: user, password: password, host: host, log: m.Log},
	}

	m.SendStatusMessage("connecting to BMC")
	successfulConnections, ecErr := common.EstablishConnections(ctx, connections)
	if ecErr != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return result, ecErr
	}
	m.SendStatusMessage("connected to BMC")

	var userAction []oob.Machine
	for _, elem := range successfulConnections {
		elem := *elem
		switch r := elem.(type) {
		case oob.Machine:
			userAction = append(userAction, r)
		}
	}
	result, err = oob.MachinePower(ctx, action, userAction)
	if err != nil {
		m.SendStatusMessage("error with " + base + ": " + err.Error())
		m.Log.V(0).Info("error with "+base, "error", err.Error())
		return result, err
	}
	m.SendStatusMessage(base + " complete")
	return result, nil
}
