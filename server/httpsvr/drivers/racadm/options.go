// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

// Options are the options racadm accepts
type Options struct {
	Address  string
	Username string
	Password string
}

// NewOptions returns an Options struct with the values provided set
func NewOptions(addr, user, pass string) Options {
	return Options{
		Address:  addr,
		Username: user,
		Password: pass,
	}
}
