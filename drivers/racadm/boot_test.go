// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

import (
	"context"
	"testing"
)

func TestBootOptions(t *testing.T) {
	s, err := NewOptions(addr, user, pass).Shell(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	dev, persistent, err := s.BootOptions()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := devices[dev]; !ok {
		t.Fatalf("unknown device: got: %s", dev)
	}
	if !persistent {
		t.Fatalf("bad peristence: want: %t, got: %t", true, persistent)
	}
}
