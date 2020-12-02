// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/bmc"
)

func init() {
	factory := func(ctx context.Context, opts bmc.DriverOptions) (bmc.Driver, error) {
		s, err := NewOptions(opts.Address, opts.Username, opts.Password).Shell(ctx)
		if err != nil {
			return nil, err
		}
		return s, nil
	}
	bmc.RegisterDriver(factory, "racadm")
}

var bmcActions = map[bmc.Action]string{
	bmc.NoAction:  "",
	bmc.ColdReset: "hard",
	bmc.WarmReset: "soft",
}

// BMC runs actions on the BMC
func (s *Shell) BMC(action bmc.Action) error {
	arg, ok := bmcActions[action]
	if !ok {
		return errors.Errorf("bmc action %q not supported by racadm driver", action)
	}
	if arg == "" {
		return nil
	}
	return s.Run("racreset" + " " + arg)
}
