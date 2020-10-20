package task

import (
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// Task interface for doing BMC actions
type Task interface {
	Execute(description string, action func(chan string) (string, oob.Error)) (id string, err error)
	Status(id string) (record repository.Record, err error)
}
