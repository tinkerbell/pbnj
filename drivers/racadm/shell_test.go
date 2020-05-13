// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

import (
	"context"
	"testing"
)

func TestShell(t *testing.T) {
	s, err := NewOptions(addr, user, pass).Shell(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	if err := s.Run("racadm serveraction powerstatus"); err != nil {
		t.Error(err)
	}
	if err := s.Close(); err != nil {
		t.Errorf("close: %s", err)
	}
}
