package rpc

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TaskService for retrieving task details.
type TaskService struct {
	TaskRunner task.Task
	v1.UnimplementedTaskServer
}

// Status returns a task record.
func (t *TaskService) Status(ctx context.Context, in *v1.StatusRequest) (*v1.StatusResponse, error) {
	l := logging.ExtractLogr(ctx)
	l.Info("start Status request", "taskID", in.TaskId)

	record, err := t.TaskRunner.Status(ctx, in.TaskId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	c := codes.OK
	if record.Error.Message != "" {
		if codes.Code(record.Error.Code) != codes.OK {
			c = codes.Code(record.Error.Code)
		} else {
			c = codes.Unknown
		}
	}

	return &v1.StatusResponse{
		Id:          record.ID,
		Description: record.Description,
		Error:       nil,
		State:       record.State,
		Result:      record.Result,
		Complete:    record.Complete,
		Messages:    record.Messages,
	}, status.Error(c, record.Error.Message)
}
