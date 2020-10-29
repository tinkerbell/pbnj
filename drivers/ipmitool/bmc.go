// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/interfaces/bmc"
)

func init() {
	factory := func(ctx context.Context, opts bmc.DriverOptions) (bmc.Driver, error) {
		s, err := NewOptions(opts.Address, opts.Username, opts.Password, opts.Cipher).Shell(ctx)
		if err != nil {
			return nil, errors.WithMessage(err, "invalid options")
		}
		return s, nil
	}
	bmc.RegisterDriver(factory, "ipmitool")
}

var bmcActions = map[bmc.Action]string{
	bmc.NoAction:        "",
	bmc.ColdReset:       "reset cold",
	bmc.WarmReset:       "reset warm",
	bmc.PassThruCommand: "command",
}

// BMC runs actions on the BMC
func (s *Shell) BMC(req bmc.BmcRequest) error {
	command, err := s.ComposeBmcCommand(req)
	if err != nil {
		return err
	}
	if command == "" {
		return nil
	}
	return s.Run(command)
}

func (s *Shell) ComposeBmcCommand(req bmc.BmcRequest) (string, error) {
	arg, ok := bmcActions[req.Action]
	if !ok {
		return "", errors.Errorf("bmc action %q not supported by ipmitool driver", req.Action)
	}
	if arg == "" {
		return "", nil
	}
	command := "bmc " + arg
	if req.Action == bmc.PassThruCommand {
		command = "bmc " + req.Command
	}
	return command, nil
}
