package rpc

import (
	"context"
	"time"

	"github.com/rs/xid"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/oob/machine"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
	"go.opentelemetry.io/otel/trace"
)

// MachineService for doing power and device actions.
type MachineService struct {
	Log logging.Logger
	// Timeout is how long a task should be run
	// before it is cancelled. This is for use in a
	// TaskRunner.Execute function that runs all BMC
	// interactions in the background.
	Timeout    time.Duration
	TaskRunner task.Task
	v1.UnimplementedMachineServer
}

// BootDevice sets the next boot device of a machine.
func (m *MachineService) BootDevice(ctx context.Context, in *v1.DeviceRequest) (*v1.DeviceResponse, error) {
	// TODO figure out how not to have to do this, but still keep the logging abstraction clean?
	l := m.Log.GetContextLogger(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)

	l.Info(
		"start BootDevice request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
		"bootDevice", in.BootDevice.String(),
		"persistent", in.Persistent,
		"efiBoot", in.EfiBoot,
	)

	execFunc := func(s chan string) (string, error) {
		mbd, err := machine.NewBootDeviceSetter(
			machine.WithDeviceRequest(in),
			machine.WithLogger(l),
			machine.WithStatusMessage(s),
		)
		if err != nil {
			return "", err
		}
		// Because this is a background task, we want to pass through the span context, but not be
		// a child context. This allows us to correctly plumb otel into the background task.
		c := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
		taskCtx, cancel := context.WithTimeout(c, m.Timeout)
		_ = cancel
		return mbd.BootDeviceSet(taskCtx, in.BootDevice.String(), in.Persistent, in.EfiBoot)
	}
	m.TaskRunner.Execute(ctx, "setting boot device", taskID, execFunc)

	return &v1.DeviceResponse{TaskId: taskID}, nil
}

// Power does a power action against a BMC.
func (m *MachineService) Power(ctx context.Context, in *v1.PowerRequest) (*v1.PowerResponse, error) {
	l := m.Log.GetContextLogger(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)
	l.Info(
		"start Power request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
		"powerAction", in.GetPowerAction().String(),
		"softTimeout", in.SoftTimeout,
		"OffDuration", in.OffDuration,
	)

	execFunc := func(s chan string) (string, error) {
		mp, err := machine.NewPowerSetter(
			machine.WithPowerRequest(in),
			machine.WithLogger(l),
			machine.WithStatusMessage(s),
		)
		if err != nil {
			return "", err
		}
		// Because this is a background task, we want to pass through the span context, but not be
		// a child context. This allows us to correctly plumb otel into the background task.
		c := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
		taskCtx, cancel := context.WithTimeout(c, m.Timeout)
		_ = cancel
		return mp.PowerSet(taskCtx, in.PowerAction.String())
	}
	m.TaskRunner.Execute(ctx, "power action: "+in.GetPowerAction().String(), taskID, execFunc)

	return &v1.PowerResponse{TaskId: taskID}, nil
}
