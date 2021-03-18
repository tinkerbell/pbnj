// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package boot

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/log"
)

var (
	logger    log.Logger
	factories = make(map[string]DriverFactory)
)

// Device describe bootable devices
type Device string

const (
	// DefaultDevice is the default boot device
	DefaultDevice = Device("default")
	// ForcePXE is used to configure pxe as the boot device
	ForcePXE = Device("pxe")
	// ForceDisk is used to configure disk as the boot device
	ForceDisk = Device("disk")
	// ForceBIOS is used to configure bios as the boot device
	ForceBIOS = Device("bios")
)

func SetupLogging(l log.Logger) {
	logger = l.Package("boot")
}

// Driver is the interface to control boot
type Driver interface {
	SetBootOptions(Options) error

	Close() error
}

// DriverFactory is used to instantiate a new Driver
type DriverFactory func(context.Context, DriverOptions) (Driver, error)

// DriverOptions contain the basic options needed to connect to device
type DriverOptions struct {
	Address  string
	Username string
	Password string
	ID       string
	Cipher   int
}

// Options contain values relevant for boot actions
type Options struct {
	// Device is the boot device to force.
	Device Device `json:"device" binding:"required"`

	// Persistent controls whether Device should apply to all future boots.
	Persistent bool `json:"persistent,omitempty"`

	// EFI controls whether the the next boot should be in EFI mode.
	EFI bool `json:"efi,omitempty"`
}

// NewDriver instantiates a new driver by calling the registered factory function
func NewDriver(ctx context.Context, name string, opts DriverOptions) (Driver, error) {
	factory, ok := factories[name]
	if !ok {
		return nil, errors.Errorf("unsupported driver type: %q", name)
	}
	return factory(ctx, opts)
}

// RegisterDriver is called by implementations in order to register their factory function for each device they can control
func RegisterDriver(factory DriverFactory, driver string) {
	if _, ok := factories[driver]; ok {
		logger.Panicf("boot driver %q already registered!", driver)
	}
	factories[driver] = factory
}
