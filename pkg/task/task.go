package task

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// Task interface for doing BMC actions.
type Task interface {
	Execute(ctx context.Context, l logr.Logger, description, taskID, host string, action func(chan string) (string, error))
	Status(ctx context.Context, taskID string) (record repository.Record, err error)
}
