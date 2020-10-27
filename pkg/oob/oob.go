package oob

import "github.com/tinkerbell/pbnj/pkg/repository"

// User management methods
type User interface {
	Create() (result string, err repository.Error)
	Update() (result string, err repository.Error)
	Delete() (result string, err repository.Error)
}

// Machine management methods
type Machine interface {
	BootDevice() (result string, err repository.Error)
	Power() (result string, err repository.Error)
}

// BMC management methods
type BMC interface {
	Reset() (result string, err repository.Error)
	NetworkSource() (result string, err repository.Error)
}
