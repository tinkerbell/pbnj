// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"os"
	"os/exec"
	"strconv"
)

var DEFAULT_CIPHER = os.Getenv("IPMITOOL_DEFAULT_CIPHER")

// ExecutablePath path to ipmitool
const ExecutablePath = "ipmitool"

// Options are the options ipmitool accepts
type Options struct {
	Address  string
	Username string
	Password string

	InterfaceName string
	Cipher        int
	Attempts      int
	RetransSecs   int
}

// NewOptions returns an Options struct with the values provided set
func NewOptions(addr, user, pass string, cipher int) Options {
	return Options{
		Address:       addr,
		Username:      user,
		Password:      pass,
		InterfaceName: "lanplus",
		Cipher:        cipher,
		Attempts:      1, // Give up quickly.
		RetransSecs:   1, // Wait 1 second between retries.
	}
}

func (o *Options) buildCommand(subcommand ...string) *exec.Cmd {
	args := make([]string, 0, len(subcommand)+6*2)

	if o.Address != "" {
		args = append(args, "-H", o.Address)
	}

	if o.Username != "" {
		args = append(args, "-U", o.Username)
	}

	if o.Password != "" {
		args = append(args, "-P", o.Password)
	}

	if o.InterfaceName != "" {
		args = append(args, "-I", o.InterfaceName)
	}

	if o.Cipher > -1 {
		args = append(args, "-C", strconv.Itoa(o.Cipher))
	} else if DEFAULT_CIPHER != "" {
		args = append(args, "-C", DEFAULT_CIPHER)
	}

	if o.Attempts > 0 {
		args = append(args, "-R", strconv.Itoa(o.Attempts))
	}

	if o.RetransSecs > 0 {
		args = append(args, "-N", strconv.Itoa(o.RetransSecs))
	}

	args = append(args, subcommand...)

	return exec.Command(ExecutablePath, args...)
}
