package machine

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/metrics"
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

// NewPowerSetter returns an oob.PowerSetter interface
func NewPowerSetter(opts ...Option) (oob.PowerSetter, error) {
	a := &Action{}
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// NewBootDeviceSetter returns an oob.BootDeviceSetter interface
func NewBootDeviceSetter(opts ...Option) (oob.BootDeviceSetter, error) {
	a := &Action{}
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// BootDeviceSet functionality for machines
func (m Action) BootDeviceSet(ctx context.Context, device string) (result string, err error) {
	labels := prometheus.Labels{
		"service": "machine",
		"action":  "boot_device",
	}
	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	host, user, password, parseErr := m.ParseAuth(m.BootDeviceRequest.Authn)
	if parseErr != nil {
		return result, parseErr
	}
	base := "setting boot device: " + m.BootDeviceRequest.GetBootDevice().String()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := map[string]interface{}{
		"ipmitool": &ipmiBootDevice{mAction: m, user: user, password: password, host: host, port: "623"},
	}

	m.SendStatusMessage("connecting to BMC")
	successfulConnections, ecErr := common.EstablishConnections(ctx, connections)
	if ecErr != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return result, ecErr
	}
	m.SendStatusMessage("connected to BMC")

	var userAction []oob.BootDeviceSetter
	for _, elem := range successfulConnections {
		conn := connections[elem]
		switch r := conn.(type) {
		case common.Connection:
			defer r.Close(ctx)
		}
		switch r := conn.(type) {
		case oob.BootDeviceSetter:
			userAction = append(userAction, r)
		}
	}
	result, err = oob.SetBootDevice(ctx, device, userAction)
	if err != nil {
		m.SendStatusMessage("error with " + base + ": " + err.Error())
		m.Log.V(0).Info("error with "+base, "error", err.Error())
		return result, err
	}
	m.SendStatusMessage(base + " complete")
	return result, nil
}

// PowerSet functionality for machines
func (m Action) PowerSet(ctx context.Context, action string) (result string, err error) {
	labels := prometheus.Labels{
		"service": "machine",
		"action":  "power",
	}
	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	host, user, password, parseErr := m.ParseAuth(m.PowerRequest.Authn)
	if parseErr != nil {
		return result, parseErr
	}
	base := "power " + m.PowerRequest.GetPowerAction().String()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	// the order here is the order in which these connections/operations will be tried

	connections := map[string]interface{}{
		"bmclib":   &bmclibBMC{user: user, password: password, host: host, log: m.Log},
		"ipmitool": &ipmiBMC{user: user, password: password, host: host, log: m.Log},
		"redfish":  &redfishBMC{user: user, password: password, host: host, log: m.Log},
	}

	successfulConnections, ecErr := common.EstablishConnections(ctx, connections)
	if ecErr != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return result, ecErr
	}
	m.SendStatusMessage("connected to BMC")

	var pwrActions []oob.PowerSetter
	for _, elem := range successfulConnections {
		conn := connections[elem]
		switch r := conn.(type) {
		case common.Connection:
			defer r.Close(ctx)
		}
		switch r := conn.(type) {
		case oob.PowerSetter:
			pwrActions = append(pwrActions, r)
		}
	}
	if len(pwrActions) == 0 {
		m.SendStatusMessage("no successful connections able to run power actions")
	}
	result, err = oob.SetPower(ctx, action, pwrActions)
	if err != nil {
		m.SendStatusMessage("error with " + base + ": " + err.Error())
		m.Log.V(0).Info("error with "+base, "error", err.Error())
		return result, err
	}
	m.SendStatusMessage(base + " complete")
	return result, nil
}
