package oob

import (
	"context"
)

// Machine management methods
type Machine interface {
	BootDevice(context.Context) (result string, err error)
	Power(context.Context) (result string, err error)
}

// BMC management methods
type BMC interface {
	// Reset() (result string, err repository.Error)
	// NetworkSource() (result string, err repository.Error)
	CreateUser(context.Context) error
	UpdateUser(context.Context) error
	DeleteUser(context.Context) error
}
