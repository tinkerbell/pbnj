package client

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "github.com/tinkerbell/pbnj/api/v1"
)

// MachinePower executes a power action against the server and retrieves status.
func MachinePower(ctx context.Context, client v1.MachineClient, taskClient v1.TaskClient, request *v1.PowerRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	response, err := client.Power(ctx, request)
	if err != nil {
		return nil, err
	}
	for to := 1; to <= 120; to++ {
		statusResp, err := taskClient.Status(ctx, &v1.StatusRequest{TaskId: response.TaskId})
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

// MachineBootDev sets the next boot device for a machine.
func MachineBootDev(ctx context.Context, client v1.MachineClient, taskClient v1.TaskClient, request *v1.DeviceRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	response, err := client.BootDevice(ctx, request)
	if err != nil {
		return nil, err
	}
	for to := 1; to <= 120; to++ {
		statusResp, err := taskClient.Status(ctx, &v1.StatusRequest{TaskId: response.TaskId})
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

// BMCCreateUser creates a BMC user.
func BMCCreateUser(ctx context.Context, client v1.BMCClient, taskClient v1.TaskClient, request *v1.CreateUserRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	response, err := client.CreateUser(ctx, request)
	if err != nil {
		return nil, err
	}
	for to := 1; to <= 120; to++ {
		statusResp, err := taskClient.Status(ctx, &v1.StatusRequest{TaskId: response.TaskId})
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

// BMCUpdateUser updates a BMC user.
func BMCUpdateUser(ctx context.Context, client v1.BMCClient, taskClient v1.TaskClient, request *v1.UpdateUserRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	response, err := client.UpdateUser(ctx, request)
	if err != nil {
		return nil, err
	}
	for to := 1; to <= 120; to++ {
		statusResp, err := taskClient.Status(ctx, &v1.StatusRequest{TaskId: response.TaskId})
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

// BMCDeleteUser updates a BMC user.
func BMCDeleteUser(ctx context.Context, client v1.BMCClient, taskClient v1.TaskClient, request *v1.DeleteUserRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	response, err := client.DeleteUser(ctx, request)
	if err != nil {
		return nil, err
	}
	for to := 1; to <= 120; to++ {
		statusResp, err := taskClient.Status(ctx, &v1.StatusRequest{TaskId: response.TaskId})
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

// Screenshot retrieves a screenshot from the server.
func Screenshot(ctx context.Context, client v1.DiagnosticClient, request *v1.ScreenshotRequest) (string, error) {
	screenshotResponse, err := client.Screenshot(ctx, request)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.%s", time.Now().String(), screenshotResponse.Filetype)

	if err := os.WriteFile(filename, screenshotResponse.Image, 0755); err != nil {
		return "", err
	}

	return filename, nil
}

// ClearSystemEventLog clears the System Event Log of the server.
func ClearSystemEventLog(ctx context.Context, client v1.DiagnosticClient, taskClient v1.TaskClient, request *v1.ClearSystemEventLogRequest) (*v1.StatusResponse, error) {
	var statusResp *v1.StatusResponse
	response, err := client.ClearSystemEventLog(ctx, request)
	if err != nil {
		return nil, err
	}

	for to := 1; to <= 120; to++ {
		statusResp, err := taskClient.Status(ctx, &v1.StatusRequest{TaskId: response.TaskId})
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

// SendNMI will tell the BMC to send an NMI to the server.
func SendNMI(ctx context.Context, client v1.DiagnosticClient, request *v1.SendNMIRequest) error {
	_, err := client.SendNMI(ctx, request)
	return err
}
