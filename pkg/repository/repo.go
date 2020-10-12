// Package repository provides primitives for
// putting tasks into a persistence layer.
package repository

import (
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
)

// Actions interface for interacting with the persistence layer
type Actions interface {
	Create(id string, val Record) error
	Get(id string) (Record, error)
	Update(id string, val Record) error
	Delete(id string) error
}

// Record that is stored in the repo
type Record struct {
	*v1.StatusResponse
}
