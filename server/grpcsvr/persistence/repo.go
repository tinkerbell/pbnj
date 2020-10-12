package persistence

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/philippgille/gokv"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// GoKV store
type GoKV struct {
	Ctx   context.Context
	Store gokv.Store
}

// Create a record
func (g *GoKV) Create(id string, val repository.Record) error {
	return g.Store.Set(id, val)
}

// Get a record
func (g *GoKV) Get(id string) (repository.Record, error) {
	rec := new(repository.Record)
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
func (g *GoKV) Update(id string, val repository.Record) error {
	rec := new(repository.Record)
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
func (g *GoKV) Delete(id string) error {
	return g.Store.Delete(id)
}
