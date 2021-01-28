package rpc

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
)

// TaskService for retrieving task details
type TaskService struct {
	Log        logging.Logger
	TaskRunner task.Task
	v1.UnimplementedTaskServer
}

// Status returns a task record
func (t *TaskService) Status(ctx context.Context, in *v1.StatusRequest) (*v1.StatusResponse, error) {
	l := t.Log.GetContextLogger(ctx)
	l.V(0).Info("start Status request", "taskID", in.TaskId)

	record, err := t.TaskRunner.Status(ctx, in.TaskId)
	if err != nil {
		return nil, err
	}
	var errMsg *v1.Error
	if record.Error != nil {
		errMsg = &v1.Error{
			Code:    record.Error.Code,
			Message: record.Error.Message,
			Details: record.Error.Details,
		}
	} else {
		errMsg = nil
	}
	// TODO should i return this as error or just return the error in the response message?
	/*if record.Error != "" {
		return nil, errors.New(record.Error)
	}*/
	return &v1.StatusResponse{
		Id:          record.Id,
		Description: record.Description,
		Error:       errMsg,
		State:       record.State,
		Result:      record.Result,
		Complete:    record.Complete,
		Messages:    record.Messages,
	}, err
}
