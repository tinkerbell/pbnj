// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

import (
	"fmt"
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
