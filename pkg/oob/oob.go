package oob

import (
	"context"

	"github.com/tinkerbell/pbnj/pkg/repository"
)

// User management methods
type User interface {
	Create() (result string, err repository.Error)
	Update() (result string, err repository.Error)
	Delete() (result string, err repository.Error)
}

// Machine management methods
type Machine interface {
	BootDevice(context.Context) (result string, err repository.Error)
	Power(context.Context) (result string, err repository.Error)
}

// BMC management methods
type BMC interface {
	// Reset() (result string, err repository.Error)
	// NetworkSource() (result string, err repository.Error)
	CreateUser(context.Context) repository.Error
	UpdateUser(context.Context) repository.Error
	DeleteUser(context.Context) repository.Error
}
