// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

import (
	"fmt"
	"github.com/tinkerbell/pbnj/interfaces/bmc"
	"os"
	"testing"
)

var (
	addr = os.Getenv("TEST_ADDR")
	user = os.Getenv("TEST_USER")
	pass = os.Getenv("TEST_PASS")
)

func TestMain(m *testing.M) {
	if addr == "" || user == "" || pass == "" {
		fmt.Println("skipping tests")
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestPassThruCommand(t *testing.T) {
	for value, req := range reqs {
		s := Shell{}
		test, err := s.ComposeBmcCommand(req)
		if err != nil {
			t.Fatal(err)
		}
		if test != value {
			t.Fatal(fmt.Sprintf("got: %s; expected: %s", test, value))
		}
	}
}

var reqs = map[string]bmc.BmcRequest{
	"racadm config -g cfgIpmiLan -o cfgIpmiLanAlertEnable 1": {
		Action:  bmc.PassThruCommand,
		Command: "config -g cfgIpmiLan -o cfgIpmiLanAlertEnable 1",
	},
	"racreset hard": {
		Action: bmc.ColdReset,
	},
	"racreset soft": {
		Action: bmc.WarmReset,
	},
}
