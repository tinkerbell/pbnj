// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"context"
	"os"
	"testing"
)

func TestShell(t *testing.T) {
	var (
		addr = os.Getenv("TEST_ADDR")
		user = os.Getenv("TEST_USER")
		pass = os.Getenv("TEST_PASS")
	)
	if addr == "" || user == "" || pass == "" {
		t.Skip()
		panic("unreachable")
	}
	s, err := NewOptions(addr, user, pass).Shell(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Run("power status"); err != nil {
		t.Error(err)
	}
	if err := s.Close(); err != nil {
		t.Errorf("close: %s", err)
	}
}
