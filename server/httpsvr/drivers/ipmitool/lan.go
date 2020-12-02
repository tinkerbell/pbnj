// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"fmt"

	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/bmc"
)

var lanIPSources = map[bmc.IPSource]string{
	bmc.IPFromDHCP: "dhcp",
	bmc.StaticIP:   "static",
}

// SetIPSource sets ip configuration method
func (s *Shell) SetIPSource(source bmc.IPSource) error {
	arg, ok := lanIPSources[source]
	if !ok {
		return fmt.Errorf("ip source %q not supported by ipmitool driver", source)
	}
	if arg == "" {
		return nil
	}
	return s.Run("lan set 1 ipsrc " + arg)
}
