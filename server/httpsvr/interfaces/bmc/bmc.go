// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package bmc

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/log"
)

var (
	logger    log.Logger
	factories = make(map[string]DriverFactory)
)

// Action is used to denote the available power actions.
type Action string

const (
	// NoAction is a NOP.
	NoAction = Action("")
	// ColdReset does a cold reset.
	ColdReset = Action("reset_cold")
	// WarmReset does a warm reset.
	WarmReset = Action("reset_warm")
)

func SetupLogging(l log.Logger) {
	logger = l.Package("bmc")
}

// ActionsBySlug is a reverse mapping of Action types.
var ActionsBySlug = map[string]Action{
	"reset_cold": ColdReset,
	"reset_warm": WarmReset,
}

// UnmarshalText unmarshals an Action from a textual representation.
func (a *Action) UnmarshalText(text []byte) error {
	if v, ok := ActionsBySlug[string(text)]; ok {
		*a = v
		return nil
	}
	return errors.Errorf("unsupported bmc action: %q", text)
}

// Driver is the interface to control BMCs.
type Driver interface {
	BMC(action Action) error
	SetIPSource(source IPSource) error

	Close() error
}

// DriverFactory is used to instantiate a new Driver.
type DriverFactory func(context.Context, DriverOptions) (Driver, error)

// DriverOptions contain the basic options needed to connect to device.
type DriverOptions struct {
	Address  string
	Username string
	Password string
	ID       string
	Cipher   int
}

// NewDriver instantiates a new driver by calling the registered factory function.
func NewDriver(ctx context.Context, name string, opts DriverOptions) (Driver, error) {
	factory, ok := factories[name]
	if !ok {
		return nil, errors.Errorf("unsupported driver type: %q", name)
	}
	return factory(ctx, opts)
}

// RegisterDriver is called by implementations in order to register their factory function for each device they can control.
func RegisterDriver(factory DriverFactory, driver string) {
	if _, ok := factories[driver]; ok {
		logger.Panicf("power driver %q already registered!", driver)
	}
	factories[driver] = factory
}
