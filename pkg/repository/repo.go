// Package repository provides primitives for
// putting tasks into a persistence layer.
package repository

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
