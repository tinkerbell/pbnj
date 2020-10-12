package grpcsvr

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/task"
	"github.com/tinkerbell/pbnj/server/grpcsvr/bmc"
)

type machineService struct {
	log        logging.Logger
	taskRunner task.Task
}

func (m *machineService) device(ctx context.Context, in *v1.DeviceRequest) (*v1.DeviceResponse, error) {
	// TODO figure out how not to have to do this, but still keep the logging abstraction clean?
	l := m.log.GetContextLogger(ctx)
	l.V(0).Info("setting boot device", "device", in.Device.String())

	switch in.GetAuthn().Authn.(type) {
	case *v1.Authn_ExternalAuthn:
		l.V(1).Info("using external authn")
	default:
		l.V(1).Info("using direct authn")
	}

	taskID, err := m.taskRunner.Execute(
		"setting boot device",
		func(s chan string) (string, oob.Error) {
			var mbd oob.Machine
			mbd = bmc.MachineAction{
				Log:               m.log,
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

func (m *machineService) powerAction(ctx context.Context, in *v1.PowerRequest) (*v1.PowerResponse, error) {
	l := m.log.GetContextLogger(ctx)
	l.V(0).Info("power request")
	// TODO INPUT VALIDATION

	switch in.GetAuthn().Authn.(type) {
	case *v1.Authn_ExternalAuthn:
		l.V(1).Info("using external authn")
	default:
		l.V(1).Info("using direct authn")
	}

	var execFunc = func(s chan string) (string, oob.Error) {
		var mp oob.Machine
		mp = bmc.MachineAction{
			Log:            m.log,
			Ctx:            ctx,
			PowerRequest:   in,
			StatusMessages: s,
		}
		return mp.Power()
	}
	taskID, err := m.taskRunner.Execute("power action", execFunc)

	return &v1.PowerResponse{
		TaskId: taskID,
	}, err
}
