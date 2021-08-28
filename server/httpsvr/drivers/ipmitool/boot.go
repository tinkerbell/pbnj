// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/boot"
)

func init() {
	factory := func(ctx context.Context, opts boot.DriverOptions) (boot.Driver, error) {
		s, err := NewOptions(opts.Address, opts.Username, opts.Password, opts.Cipher).Shell(ctx)
		if err != nil {
			return nil, errors.WithMessage(err, "invalid options")
		}
		return s, nil
	}
	boot.RegisterDriver(factory, "ipmitool")
}

var devices = map[boot.Device]string{
	boot.DefaultDevice: "none",
	boot.ForcePXE:      "pxe",
	boot.ForceDisk:     "disk",
	boot.ForceBIOS:     "bios",
}

// SetBootOptions sets boot options.
func (s *Shell) SetBootOptions(opts boot.Options) error {
	device, ok := devices[opts.Device]
	if !ok {
		// TODO(betawaffle): Make this a better error type.
		return errors.Errorf("boot device %q not supported by ipmitool driver", opts.Device)
	}
	if device == "" {
		return nil
	}
	var options []string
	if opts.Persistent {
		options = append(options, "persistent")
	}
	if opts.EFI {
		options = append(options, "efiboot")
	}

	cmd := fmt.Sprintf("chassis bootdev %s", device)
	if len(options) > 0 {
		cmd = fmt.Sprintf("%s options=%s", cmd, strings.Join(options, ","))
	}
	return s.Run(cmd)
}
