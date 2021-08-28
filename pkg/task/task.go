package task

import (
	"context"

	"github.com/tinkerbell/pbnj/pkg/repository"
)

// Task interface for doing BMC actions.
type Task interface {
	Execute(ctx context.Context, description, taskID string, action func(chan string) (string, error))
	Status(ctx context.Context, taskID string) (record repository.Record, err error)
}
