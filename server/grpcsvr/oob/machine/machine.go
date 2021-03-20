package machine

import (
	"context"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-logr/logr"
	"github.com/jacobweinstock/registrar"
	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/repository"
	common "github.com/tinkerbell/pbnj/server/grpcsvr/oob"
)

// Action for making power actions on BMCs, implements oob.Machine interface
type Action struct {
	common.Accessory
	PowerRequest      *v1.PowerRequest
	BootDeviceRequest *v1.DeviceRequest
}

// Option to add to an Actions
type Option func(a *Action)

// WithLogger adds a logr to an Action struct
func WithLogger(l logr.Logger) Option {
	return func(a *Action) { a.Log = l }
}

// WithStatusMessage adds a status message chan to an Action struct
func WithStatusMessage(s chan string) Option {
	return func(a *Action) { a.StatusMessages = s }
}

// WithDeviceRequest adds DeviceRequest to an Action struct
func WithDeviceRequest(in *v1.DeviceRequest) Option {
	return func(a *Action) { a.BootDeviceRequest = in }
}

// WithPowerRequest adds PowerRequest to an Action struct
func WithPowerRequest(in *v1.PowerRequest) Option {
	return func(a *Action) { a.PowerRequest = in }
}

// NewPowerSetter returns an oob.PowerSetter interface
func NewPowerSetter(opts ...Option) *Action {
	a := &Action{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// NewBootDeviceSetter returns an oob.BootDeviceSetter interface
func NewBootDeviceSetter(opts ...Option) *Action {
	a := &Action{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// BootDeviceSet functionality for machines
func (m Action) BootDeviceSet(ctx context.Context, device string, persistent, efiBoot bool) (result string, err error) {
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
	client := bmclib.NewClient(host, "623", user, password, bmclib.WithLogger(m.Log))
	//client.Registry.Drivers = client.Registry.FilterForCompatible(ctx)

	m.SendStatusMessage("connecting to BMC")
	err = client.Open(ctx)
	if err != nil {
		return "", &repository.Error{
			Code:    v1.Code_value["PERMISSION_DENIED"],
			Message: err.Error(),
		}
	}
	defer client.Close(ctx)

	var dev string
	switch device {
	case v1.BootDevice_BOOT_DEVICE_NONE.String():
		dev = "none"
	case v1.BootDevice_BOOT_DEVICE_BIOS.String():
		dev = "bios"
	case v1.BootDevice_BOOT_DEVICE_CDROM.String():
		dev = "cdrom"
	case v1.BootDevice_BOOT_DEVICE_DISK.String():
		dev = "disk"
	case v1.BootDevice_BOOT_DEVICE_PXE.String():
		dev = "pxe"
	case v1.BootDevice_BOOT_DEVICE_UNSPECIFIED.String():
		return "", &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "UNSPECIFIED boot device",
		}
	default:
		return "", &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "unknown boot device",
		}
	}

	var errMsg string
	ok, err := client.SetBootDevice(ctx, dev, persistent, efiBoot)
	if err != nil {
		errMsg = err.Error()
	} else if !ok {
		errMsg = "setting boot device failed"
	}
	if errMsg != "" {
		return "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: errMsg,
		}
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

	client := bmclib.NewClient(host, "623", user, password, bmclib.WithLogger(m.Log))
	client.Registry.Register("pureipmi", "ipmi", registrar.Features{""}, nil, &gebnConn{user: user, password: password, host: host, log: m.Log})
	client.Registry.Register("gofish", "redfish", registrar.Features{""}, nil, &redfishConn{user: user, password: password, host: host, log: m.Log})
	err = client.Open(ctx)
	if err != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return result, &repository.Error{
			Code:    v1.Code_value["PERMISSION_DENIED"],
			Message: err.Error(),
		}
	}
	defer client.Close(ctx)
	m.SendStatusMessage("connected to BMC")

	var ok bool
	if action == v1.PowerAction_POWER_ACTION_STATUS.String() {
		result, err = client.GetPowerState(ctx)
		ok = true
	} else {
		switch action {
		case v1.PowerAction_POWER_ACTION_ON.String():
			result = "on"
			ok, err = client.SetPowerState(ctx, result)
		case v1.PowerAction_POWER_ACTION_OFF.String():
			result = "soft"
			ok, err = client.SetPowerState(ctx, result)
		case v1.PowerAction_POWER_ACTION_RESET.String():
			result = "reset"
			ok, err = client.SetPowerState(ctx, result)
		case v1.PowerAction_POWER_ACTION_HARDOFF.String():
			result = "off"
			ok, err = client.SetPowerState(ctx, result)
		case v1.PowerAction_POWER_ACTION_CYCLE.String():
			result = "cycle"
			ok, err = client.SetPowerState(ctx, result)
		case v1.PowerAction_POWER_ACTION_UNSPECIFIED.String():
			return "", &repository.Error{
				Code:    v1.Code_value["INVALID_ARGUMENT"],
				Message: "UNSPECIFIED power action",
			}
		default:
			return "", &repository.Error{
				Code:    v1.Code_value["INVALID_ARGUMENT"],
				Message: "unknown power action",
			}
		}
	}
	if err != nil {
		m.SendStatusMessage("error with " + base + ": " + err.Error())
		return "", &repository.Error{
			Code:    v1.Code_value["INTERNAL"],
			Message: "failed power action",
		}
	}
	if !ok {
		m.SendStatusMessage("problem with " + base)
		return "", &repository.Error{
			Code:    v1.Code_value["INTERNAL"],
			Message: "failed power action",
		}
	}
	m.SendStatusMessage(base + " complete")
	return result, nil
}
