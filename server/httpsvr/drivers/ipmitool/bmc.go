// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/bmc"
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
	bmc.NoAction:  "",
	bmc.ColdReset: "reset cold",
	bmc.WarmReset: "reset warm",
}

// BMC runs actions on the BMC
func (s *Shell) BMC(action bmc.Action) error {
	arg, ok := bmcActions[action]
	if !ok {
		return errors.Errorf("bmc action %q not supported by ipmitool driver", action)
	}
	if arg == "" {
		return nil
	}
	return s.Run("bmc " + arg)
}
