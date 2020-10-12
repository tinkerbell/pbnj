package client

import (
	"context"
	"time"

	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
)

// MachinePower executes a power action against the server and retrieves status
func MachinePower(client v1.MachineClient, taskClient v1.TaskClient, request *v1.PowerRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
