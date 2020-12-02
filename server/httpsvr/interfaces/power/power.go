// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package power

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/evlog"
	"github.com/tinkerbell/pbnj/server/httpsvr/log"
)

var (
	logger    log.Logger
	elog      *evlog.Log
	factories = make(map[string]DriverFactory)
)

// Action is the type for available actions
type Action string

func SetupLogging(l log.Logger) {
	logger = l.Package("power")
	elog = evlog.New(logger)
}

const (
	// NoAction is a Nop
	NoAction = Action("")
	// TurnOn will power up the device
	TurnOn = Action("turn_on")
	// SoftOff will perform a soft-shutdown (usually through OS via ACPI)
	SoftOff = Action("soft_off")
	// HardOff will power off the device without informing the OS
	HardOff = Action("hard_off")
	// Reset will perform a hard reset without informing the OS
	Reset = Action("reset")
)

// Driver is the interface to control power
type Driver interface {
	PowerStatus() (Status, error)
	Power(action Action) error

	LastStatus() Status

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

// Options contain time values relevant for power actions
type Options struct {
	SoftTimeout    time.Duration
	OffTimeout     time.Duration
	OffDuration    time.Duration
	OnTimeout      time.Duration
	IgnoreRunError bool
}

// DefaultOptions are the default times for power actions
var DefaultOptions = Options{
	SoftTimeout:    30 * time.Second,
	OffDuration:    3 * time.Second,
	OffTimeout:     30 * time.Second,
	OnTimeout:      30 * time.Second,
	IgnoreRunError: false,
}

// Status is the type for reporting power status
type Status string

const (
	// AnyStatus is a fallback for unknown status
	AnyStatus = Status("")
	// Off is returned when the device is powered off
	Off = Status("off")
	// On is returned when the device is powered on
	On = Status("on")
)

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
		logger.Panicf("power driver %q already registered!", driver)
	}
	factories[driver] = factory
}
