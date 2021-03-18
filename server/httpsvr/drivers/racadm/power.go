// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

import (
	"context"
	"regexp"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/power"
)

var powerStatusRegexp = regexp.MustCompile(`(?m)^1\r?\n`)

func init() {
	factory := func(ctx context.Context, opts power.DriverOptions) (power.Driver, error) {
		s, err := NewOptions(opts.Address, opts.Username, opts.Password).Shell(ctx)
		if err != nil {
			return nil, errors.WithMessage(err, "initializing racadm power")
		}
		return s, nil
	}
	power.RegisterDriver(factory, "racadm")
}

var powerActions = map[power.Action]string{
	power.NoAction: "",
	power.TurnOn:   "powerup",
	power.SoftOff:  "graceshutdown",
	power.HardOff:  "powerdown",
	power.Reset:    "hardreset",
}

// Power sets the power state
func (s *Shell) Power(action power.Action) error {
	if action == power.NoAction {
		return nil
	}

	arg, ok := powerActions[action]
	if !ok {
		// TODO(betawaffle): Make this a better error type.
		return errors.Errorf("power action %q not supported by racadm driver", action)
	}

	return s.Run("racadm serveraction" + " " + arg)
}

// PowerStatus returns the power state
func (s *Shell) PowerStatus() (power.Status, error) {
	out, err := s.Output("racadm get system.power.status")
	if err != nil {
		return power.AnyStatus, errors.WithMessage(err, "failure getting system.power.status")
	}

	status := power.Off

	if powerStatusRegexp.MatchString(out) {
		status = power.On
	}

	s.lock.Lock()
	s.lastStatus = status
	s.lock.Unlock()

	return status, err
}

// LastStatus returns the previous power status
func (s *Shell) LastStatus() power.Status {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.lastStatus
}
