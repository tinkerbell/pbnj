package machine

import (
	"context"
	"fmt"
	"strings"

	"github.com/bmc-toolbox/bmclib"
	"github.com/bmc-toolbox/bmclib/bmc"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/repository"
	common "github.com/tinkerbell/pbnj/server/grpcsvr/oob"
)

// Action for making power actions on BMCs, implements oob.Machine interface.
type Action struct {
	common.Accessory
	PowerRequest      *v1.PowerRequest
	BootDeviceRequest *v1.DeviceRequest
}

// Option to add to an Actions.
type Option func(a *Action) error

// WithLogger adds a logr to an Action struct.
func WithLogger(l logr.Logger) Option {
	return func(a *Action) error {
		a.Log = l
		return nil
	}
}

// WithStatusMessage adds a status message chan to an Action struct.
func WithStatusMessage(s chan string) Option {
	return func(a *Action) error {
		a.StatusMessages = s
		return nil
	}
}

// WithDeviceRequest adds DeviceRequest to an Action struct.
func WithDeviceRequest(in *v1.DeviceRequest) Option {
	return func(a *Action) error {
		a.BootDeviceRequest = in
		return nil
	}
}

// WithPowerRequest adds PowerRequest to an Action struct.
func WithPowerRequest(in *v1.PowerRequest) Option {
	return func(a *Action) error {
		a.PowerRequest = in
		return nil
	}
}

// NewPowerSetter returns an oob.PowerSetter interface.
func NewPowerSetter(opts ...Option) (*Action, error) {
	a := &Action{}
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// NewBootDeviceSetter returns an oob.BootDeviceSetter interface.
func NewBootDeviceSetter(opts ...Option) (*Action, error) {
	a := &Action{}
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// BootDeviceSet functionality for machines.
func (m Action) BootDeviceSet(ctx context.Context, device string, persistent, efiBoot bool) (result string, err error) {
	labels := prometheus.Labels{
		"service": "machine",
		"action":  "boot_device",
	}
	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

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

	host, user, password, parseErr := m.ParseAuth(m.BootDeviceRequest.Authn)
	if parseErr != nil {
		return result, parseErr
	}
	base := "setting boot device: " + m.BootDeviceRequest.GetBootDevice().String()
	msg := "working on " + base
	m.SendStatusMessage(msg)
	client := bmclib.NewClient(host, "623", user, password, bmclib.WithLogger(m.Log))

	m.SendStatusMessage("connecting to BMC")
	err = client.Open(ctx)
	if err != nil {
		return "", &repository.Error{
			Code:    v1.Code_value["PERMISSION_DENIED"],
			Message: err.Error(),
		}
	}
	log := m.Log.WithValues("device", dev, "host", host, "user", user)
	defer func() {
		client.Close(ctx)
		log.Info("closed connections", logMetadata(client.GetMetadata())...)
	}()
	log.Info("connected to BMC", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage("connected to BMC")

	ok, err := client.SetBootDevice(ctx, dev, persistent, efiBoot)
	log = m.Log.WithValues(logMetadata(client.GetMetadata()))
	if err != nil {
		log.Error(err, "failed to set boot device", "device", dev)
	} else if !ok {
		err = fmt.Errorf("setting boot device failed")
	}
	if err != nil {
		log.Error(err, "error with "+base)
		m.SendStatusMessage(fmt.Sprintf("failed to set %v as boot device", dev))

		return "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	log.Info(base + " complete")
	m.SendStatusMessage(base + " complete")

	return result, nil
}

// PowerSet functionality for machines.
func (m Action) PowerSet(ctx context.Context, action string) (result string, err error) {
	labels := prometheus.Labels{
		"service": "machine",
		"action":  "power",
	}
	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	var pwrAction string
	switch action {
	case v1.PowerAction_POWER_ACTION_ON.String():
		pwrAction = "on"
	case v1.PowerAction_POWER_ACTION_OFF.String():
		pwrAction = "off"
	case v1.PowerAction_POWER_ACTION_STATUS.String():
		pwrAction = "status"
	case v1.PowerAction_POWER_ACTION_RESET.String():
		pwrAction = "reset"
	case v1.PowerAction_POWER_ACTION_HARDOFF.String():
		pwrAction = "off"
	case v1.PowerAction_POWER_ACTION_CYCLE.String():
		pwrAction = "cycle"
	case v1.PowerAction_POWER_ACTION_UNSPECIFIED.String():
		return "", &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "UNSPECIFIED power action",
		}
	default:
		return "", &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: fmt.Sprintf("unknown power action: %q", action),
		}
	}

	host, user, password, parseErr := m.ParseAuth(m.PowerRequest.Authn)
	if parseErr != nil {
		return result, parseErr
	}
	base := "power " + m.PowerRequest.GetPowerAction().String()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	client := bmclib.NewClient(host, "623", user, password, bmclib.WithLogger(m.Log))

	err = client.Open(ctx)
	if err != nil {
		m.SendStatusMessage("connecting to BMC failed")

		return "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	log := m.Log.WithValues("action", action, "host", host, "user", user)

	defer func() {
		client.Close(ctx)
		log.Info("closed connections", logMetadata(client.GetMetadata())...)
	}()
	log.Info("connected to BMC", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage("connected to BMC")

	ok := true
	if pwrAction == "status" {
		result, err = client.GetPowerState(ctx)
	} else {
		if action == v1.PowerAction_POWER_ACTION_CYCLE.String() {
			// check status
			// if powered on, do cycle
			// if powered off, do power on
			status, err := client.GetPowerState(ctx)
			if err != nil {
				log.V(0).Error(err, "failed to set power state "+base)
				m.SendStatusMessage("error with " + base + ": " + err.Error())
				return "", &repository.Error{
					Code:    v1.Code_value["UNKNOWN"],
					Message: err.Error(),
				}
			}
			if strings.Contains(strings.ToLower(status), "off") {
				pwrAction = v1.PowerAction_POWER_ACTION_ON.String()
			}
		}
		ok, err = client.SetPowerState(ctx, pwrAction)
		result = fmt.Sprintf("%v complete", base)
	}
	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	if err != nil {
		log.Error(err, "failed to set power state "+base)
		m.SendStatusMessage("error with " + base + ": " + err.Error())
	}
	if !ok && err == nil {
		log.Error(err, fmt.Sprintf("error completing power %v action", action))
		err = fmt.Errorf("error completing power %v action", action)
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("error completing power %v action", action))

		return "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	log.Info(base + " complete")
	m.SendStatusMessage(base + " complete")

	return result, nil
}

func logMetadata(md bmc.Metadata) []interface{} {
	kvs := []interface{}{
		"ProvidersAttempted", md.ProvidersAttempted,
		"SuccessfulOpenConns", md.SuccessfulOpenConns,
		"SuccessfulCloseConns", md.SuccessfulCloseConns,
		"SuccessfulProvider", md.SuccessfulProvider,
	}

	return kvs
}
