// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/interfaces/power"
)

func init() {
	factory := func(ctx context.Context, opts power.DriverOptions) (power.Driver, error) {
		s, err := NewOptions(opts.Address, opts.Username, opts.Password).Shell(ctx)
		if err != nil {
			return nil, errors.WithMessage(err, "invalid options")
		}
		return s, nil
	}
	power.RegisterDriver(factory, "ipmitool")
}

var powerActions = map[power.Action]string{
	power.NoAction: "",
	power.TurnOn:   "on",
	power.SoftOff:  "soft",
	power.HardOff:  "off",
	power.Reset:    "reset",
}

// Power sets the power state
func (s *Shell) Power(action power.Action) error {
	arg, ok := powerActions[action]
	if !ok {
		// TODO(betawaffle): Make this a better error type.
		return errors.Errorf("power action %q not supported by ipmitool driver", action)
	}
	if arg == "" {
		return nil
	}

	return s.Run("power " + arg)
}

// PowerStatus returns the power state
func (s *Shell) PowerStatus() (power.Status, error) {
	err := s.Run("power status")
	return s.LastStatus(), errors.WithMessage(err, "error retrieving ipmi power status")
}
