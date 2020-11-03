package client

import (
	"context"
	"time"

	v1 "github.com/tinkerbell/pbnj/api/v1"
)

// MachinePower executes a power action against the server and retrieves status
func MachinePower(ctx context.Context, client v1.MachineClient, taskClient v1.TaskClient, request *v1.PowerRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	response, err := client.Power(ctx, request)
	if err != nil {
		return nil, err
	}
	for to := 1; to <= 120; to++ {
		statusResp, err := taskClient.Status(ctx, &v1.StatusRequest{
			TaskId: response.TaskId,
		})
		if err != nil {
			return nil, err
		}
		if statusResp.Complete {
			return statusResp, nil
		}
		time.Sleep(1 * time.Second)
	}
	return statusResp, nil

}

// MachineBootDev sets the next boot device for a machine
func MachineBootDev(ctx context.Context, client v1.MachineClient, taskClient v1.TaskClient, request *v1.DeviceRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	response, err := client.BootDevice(ctx, request)
	if err != nil {
		return nil, err
	}
	for to := 1; to <= 120; to++ {
		statusResp, err := taskClient.Status(ctx, &v1.StatusRequest{
			TaskId: response.TaskId,
		})
		if err != nil {
			return nil, err
		}
		if statusResp.Complete {
			return statusResp, nil
		}
		time.Sleep(1 * time.Second)
	}
	return statusResp, nil
}
