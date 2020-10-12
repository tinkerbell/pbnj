package oob

import (
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
)

// User management methods
type User interface {
	Create() (result string, err *Error)
	Update() (result string, err *Error)
	Delete() (result string, err *Error)
}

// Machine management methods
type Machine interface {
	BootDevice() (result string, err Error)
	Power() (result string, err Error)
}

// BMC management methods
type BMC interface {
	Reset() (result string, err Error)
	NetworkSource() (result string, err Error)
}

// Error for all bmc actions
type Error struct {
	v1.Error
}
