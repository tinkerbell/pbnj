package machine

import (
	"context"
	"fmt"
	"strings"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/oob"
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

// NewAction returns a new Action
func NewAction(opts ...Option) (*Action, error) {
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
	two := &bmclibBMC{user: user, password: password, host: host, log: m.Log}
	client.Registry.Register("bmclib.legacy", "", nil, "", two)
	client.Registry.Drivers = client.Registry.For("bmclib.legacy")
	err = client.Open(ctx)
	if err != nil {
		return "", &repository.Error{
			Code:    v1.Code_value["PERMISSION_DENIED"],
			Message: err.Error(),
		}
	}
	defer client.Close(ctx)
	m.SendStatusMessage("connected to BMC")

	var ok bool
	var errMsg string
	if action == v1.PowerAction_POWER_ACTION_STATUS.String() {
		result, err = client.GetPowerState(ctx)
	} else {
		ok, err = client.SetPowerState(ctx, lookupPower(action))
		if !ok {
			errMsg = "power set failed"
		} else {
			result = action
		}
	}

	if err != nil {
		errMsg = err.Error()
	}
	if errMsg != "" {
		err := &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: errMsg,
		}
		return "", err
	}
	s := strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%+v\n", client.GetMetadata()), "}", ""), "{", "")
	s = strings.TrimSuffix(s, "\n")
	m.SendStatusMessage(s)
	m.SendStatusMessage(base + " complete")
	return result, nil
}

func lookupPower(p string) string {
	var o string
	switch p {
	case v1.PowerAction_POWER_ACTION_ON.String():
		o = "on"
	case v1.PowerAction_POWER_ACTION_HARDOFF.String():
		o = "off"
	case v1.PowerAction_POWER_ACTION_OFF.String():
		o = "soft"
	case v1.PowerAction_POWER_ACTION_RESET.String():
		o = "reset"
	case v1.PowerAction_POWER_ACTION_CYCLE.String():
		o = "cycle"
	case v1.PowerAction_POWER_ACTION_UNSPECIFIED.String():
		o = "unknown"
	default:
		o = "unknown"
	}
	return o
}
