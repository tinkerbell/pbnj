package rpc

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
)

// TaskService for retrieving task details
type TaskService struct {
	Log        logging.Logger
	TaskRunner task.Task
}

// Task returns a task record
func (t *TaskService) Task(ctx context.Context, in *v1.StatusRequest) (*v1.StatusResponse, error) {
	l := t.Log.GetContextLogger(ctx)
	l.V(0).Info("getting task record")

	record, err := t.TaskRunner.Status(in.TaskId)
	if err != nil {
		l.V(0).Error(err, "error getting task status")
		return nil, err
	}
	// TODO should i return this as error or just return the error in the response message?
	/*if record.Error != "" {
		return nil, errors.New(record.Error)
	}*/
	return &v1.StatusResponse{
		Id:          record.Id,
		Description: record.Description,
		Error: &v1.Error{
			Code:    record.Error.Code,
			Message: record.Error.Message,
			Details: record.Error.Details,
		},
		State:    record.State,
		Result:   record.Result,
		Complete: record.Complete,
		Messages: record.Messages,
	}, err
}
