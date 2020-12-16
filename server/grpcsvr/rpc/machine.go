package rpc

import (
	"context"

	"github.com/pkg/errors"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
	"github.com/tinkerbell/pbnj/server/grpcsvr/oob/machine"
)

// MachineService for doing power and device actions
type MachineService struct {
	Log        logging.Logger
	TaskRunner task.Task
	v1.UnimplementedMachineServer
}

// BootDevice sets the next boot device of a machine
func (m *MachineService) BootDevice(ctx context.Context, in *v1.DeviceRequest) (*v1.DeviceResponse, error) {
	// TODO figure out how not to have to do this, but still keep the logging abstraction clean?
	l := m.Log.GetContextLogger(ctx)
	l.V(0).Info("setting boot device", "device", in.BootDevice.String())
	if err := in.Validate(); err != nil {
		return nil, errors.Wrap(err, "input arguments are invalid")
	}

	taskID, err := m.TaskRunner.Execute(
		ctx,
		"setting boot device",
		func(s chan string) (string, error) {
			mbd, err := machine.NewBootDeviceSetter(
				machine.WithDeviceRequest(in),
				machine.WithLogger(l),
				machine.WithStatusMessage(s),
			)
			if err != nil {
				return "", err
			}
			taskCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			_ = cancel
			return mbd.BootDeviceSet(taskCtx, in.BootDevice.String())
		})

	return &v1.DeviceResponse{
		TaskId: taskID,
	}, err
}

// Power does a power action against a BMC
func (m *MachineService) Power(ctx context.Context, in *v1.PowerRequest) (*v1.PowerResponse, error) {
	l := m.Log.GetContextLogger(ctx)
	l.V(0).Info("power request")
	if err := in.Validate(); err != nil {
		return nil, errors.Wrap(err, "input arguments are invalid")
	}

	var execFunc = func(s chan string) (string, error) {
		mp, err := machine.NewPowerSetter(
			machine.WithPowerRequest(in),
			machine.WithLogger(l),
			machine.WithStatusMessage(s),
		)
		if err != nil {
			return "", err
		}
		taskCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		_ = cancel
		return mp.PowerSet(taskCtx, in.PowerAction.String())
	}
	taskID, err := m.TaskRunner.Execute(ctx, "power action: "+in.GetPowerAction().String(), execFunc)

	return &v1.PowerResponse{
		TaskId: taskID,
	}, err
}
