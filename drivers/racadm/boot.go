// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/interfaces/boot"
)

func init() {
	factory := func(ctx context.Context, opts boot.DriverOptions) (boot.Driver, error) {
		s, err := NewOptions(opts.Address, opts.Username, opts.Password).Shell(ctx)
		if err != nil {
			return nil, errors.WithMessage(err, "initializing racadm boot")
		}
		return s, nil
	}
	boot.RegisterDriver(factory, "racadm")
}

var devices = map[boot.Device]string{
	boot.DefaultDevice: "normal",
	boot.ForcePXE:      "pxe",
	boot.ForceDisk:     "hdd",
	boot.ForceBIOS:     "bios",
}

var dell2Devices = map[string]boot.Device{
	"Normal": boot.DefaultDevice,
	"PXE":    boot.ForcePXE,
	"HDD":    boot.ForceDisk,
	"BIOS":   boot.ForceBIOS,
}

// SetBootOptions sets boot options
func (s *Shell) SetBootOptions(opts boot.Options) error {
	device, ok := devices[opts.Device]
	if !ok {
		// TODO(betawaffle): Make this a better error type.
		return errors.Errorf("boot device %q not supported by racadm driver", opts.Device)
	}

	bootonce := "1"
	if opts.Persistent {
		bootonce = "0"
	}

	cmd := "racadm set idrac.serverboot.firstbootdevice" + " " + device
	err := s.Run(cmd)
	if err != nil {
		return errors.WithMessage(err, "failure setting idrac.serverboot.firstbootdevice")
	}
	cmd = "racadm set idrac.serverboot.bootonce" + " " + bootonce
	return s.Run(cmd)
}

// BootOptions returns the configured boot options
func (s *Shell) BootOptions() (boot.Device, bool, error) {
	devLine, err := s.Output("racadm get idrac.serverboot.firstbootdevice")
	if err != nil {
		return "", false, errors.WithMessage(err, "failure getting idrac.serverboot.firstbootdevice")
	}
	dellDev := strings.ToLower(strings.Split(devLine, "=")[1])
	device, ok := dell2Devices[dellDev]
	if !ok {
		return "", false, errors.New("unknown racadm device, device=\"" + dellDev + "\"")
	}

	persistent, err := s.Output("racadm get idrac.serverboot.bootonce")
	if err != nil {
		return "", false, errors.WithMessage(err, "failure getting idrac.serverboot.bootonce")
	}

	return device, persistent == "BootOnce=Enabled", nil
}
