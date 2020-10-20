// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package evlog

import (
	"fmt"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	tx := Tx{ID: "Tx"}
	trace := &Trace{
		start:  time.Now(),
		event:  "event",
		fields: []interface{}{"field1", "1", "field2", "2"},
		tx:     tx,
	}
	got := trace.String()
	want := "txid=Tx event=event_start field1=1 field2=2"
	if got != want {
		t.Fatalf("want: |%s|, got: |%s|", want, got)
	}

	trace.err = fmt.Errorf("some error")
	got = trace.String()
	want = "txid=Tx event=event_start err=some error field1=1 field2=2"
	if got != want {
		t.Fatalf("want: |%s|, got: |%s|", want, got)
	}

	trace.end = time.Now()
	got = trace.String()
	duration := trace.end.Sub(trace.start).String()
	want = "txid=Tx duration=" + duration + " event=event_end err=some error field1=1 field2=2"
	if got != want {
		t.Fatalf("want: |%s|, got: |%s|", want, got)
	}
}
