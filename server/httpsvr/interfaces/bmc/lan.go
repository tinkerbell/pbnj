// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package bmc

import "github.com/pkg/errors"

// IPSource denotes the types of ip configurations
type IPSource string

const (
	// IPFromDHCP denotes ip via dhcp configuration
	IPFromDHCP = IPSource("dhcp")
	// StaticIP denotes ip via static configuration
	StaticIP = IPSource("static")
)

// IPSourceBySlug is a reverse mapping of strings to IPSource
var IPSourceBySlug = map[string]IPSource{
	"dhcp":   IPFromDHCP,
	"static": StaticIP,
}

// UnmarshalText unmarshals an IPSource from a textual representation
func (s *IPSource) UnmarshalText(text []byte) error {
	if v, ok := IPSourceBySlug[string(text)]; ok {
		*s = v
		return nil
	}
	return errors.Errorf("unsupported bmc lan ip source: %q", text)
}
