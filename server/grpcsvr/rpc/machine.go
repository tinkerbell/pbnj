package rpc

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/repository"
	"github.com/tinkerbell/pbnj/pkg/task"
	"github.com/tinkerbell/pbnj/server/grpcsvr/bmc"
)

// MachineService for doing power and device actions
type MachineService struct {
	Log        logging.Logger
	TaskRunner task.Task
}

// Device sets the next boot device of a machine
func (m *MachineService) Device(ctx context.Context, in *v1.DeviceRequest) (*v1.DeviceResponse, error) {
	// TODO figure out how not to have to do this, but still keep the logging abstraction clean?
	l := m.Log.GetContextLogger(ctx)
	l.V(0).Info("setting boot device", "device", in.Device.String())

	taskID, err := m.TaskRunner.Execute(
		"setting boot device",
		func(s chan string) (string, repository.Error) {
			mbd := bmc.MachineAction{
				Log:               m.Log,
				Ctx:               ctx,
				BootDeviceRequest: in,
				StatusMessages:    s,
			}
			return mbd.BootDevice()
		})

	return &v1.DeviceResponse{
		TaskId: taskID,
	}, err
}

// PowerAction does a power action against a BMC
func (m *MachineService) PowerAction(ctx context.Context, in *v1.PowerRequest) (*v1.PowerResponse, error) {
	l := m.Log.GetContextLogger(ctx)
	l.V(0).Info("power request")
	// TODO INPUT VALIDATION

	var execFunc = func(s chan string) (string, repository.Error) {
		mp := bmc.MachineAction{
			Log:            m.Log,
			Ctx:            ctx,
			PowerRequest:   in,
			StatusMessages: s,
		}
		return mp.Power()
	}
	taskID, err := m.TaskRunner.Execute("power action: "+in.GetAction().String(), execFunc)

	return &v1.PowerResponse{
		TaskId: taskID,
	}, err
}
