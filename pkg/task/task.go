package task

import (
	"context"

	"github.com/tinkerbell/pbnj/pkg/repository"
)

// Task interface for doing BMC actions
type Task interface {
	Execute(ctx context.Context, description string, action func(chan string) (string, error)) (id string, err error)
	Status(ctx context.Context, id string) (record repository.Record, err error)
}
