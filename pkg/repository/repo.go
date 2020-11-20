// Package repository provides primitives for
// putting tasks into a persistence layer.
package repository

import "fmt"

// Actions interface for interacting with the persistence layer
type Actions interface {
	Create(id string, val Record) error
	Get(id string) (Record, error)
	Update(id string, val Record) error
	Delete(id string) error
}

// Record that is stored in the repo
type Record struct {
	//*v1.StatusResponse
	Id          string
	Description string
	Error       *Error
	State       string
	Result      string
	Complete    bool
	Messages    []string
}

// Error for all bmc actions
type Error struct {
	Code    int32
	Message string
	Details []string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %v message: %v details: %v", e.Code, e.Message, e.Details)
}

// StructuredError returns the error struct for convenience
func (e *Error) StructuredError() *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	}
}
