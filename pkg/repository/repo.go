// Package repository provides primitives for
// putting tasks into a persistence layer.
package repository

import (
	"context"
	"fmt"

	"github.com/philippgille/gokv"
	"github.com/pkg/errors"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
)

// Actions interface for interacting with the persistence layer
type Actions interface {
	Create(ctx context.Context, id string, val Record) error
	Get(ctx context.Context, id string) (Record, error)
	Update(ctx context.Context, id string, val Record) error
	Delete(ctx context.Context, id string) error
}

// Record that is stored in the repo
type Record struct {
	*v1.StatusResponse
}

// GoKV store
type GoKV struct {
	Store gokv.Store
}

// Create a record
func (g *GoKV) Create(ctx context.Context, id string, val Record) error {
	return g.Store.Set(id, val)
}

// Get a record
func (g *GoKV) Get(ctx context.Context, id string) (Record, error) {
	rec := new(Record)
	found, err := g.Store.Get(id, rec)
	if err != nil {
		return *rec, err
	}
	if !found {
		err = errors.New(fmt.Sprintf("record id not found: %v", id))
	}
	return *rec, err
}

// Update a record
func (g *GoKV) Update(ctx context.Context, id string, val Record) error {
	rec := new(Record)
	found, err := g.Store.Get(id, rec)
	if err != nil {
		return err
	}
	if !found {
		return errors.New(fmt.Sprintf("record id not found: %v", id))
	}
	return g.Store.Set(id, val)
}

// Delete a record
func (g *GoKV) Delete(ctx context.Context, id string) error {
	return g.Store.Delete(id)
}
