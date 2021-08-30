// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package evlog

import (
	"strings"
	"time"
)

// Trace is used for detailed tracing of events.
type Trace struct {
	tx Tx

	event  string
	start  time.Time
	end    time.Time
	err    error
	fields []interface{}
}

// Stop records the time, along with any non-nil error.
func (t *Trace) Stop(err *error) {
	t.end = time.Now()

	if err != nil {
		t.err = *err
	}

	t.tx.l.logger.With(t.tx.getFields()...).With(t.getFields()...).Debug()
}

func (t *Trace) getFields() []interface{} {
	event := t.event

	fields := make([]interface{}, 0, 6+len(t.fields))
	switch {
	case !t.end.IsZero():
		event += "_end"
		fields = append(fields, "duration", t.end.Sub(t.start))
	case !t.start.IsZero():
		event += "_start"
	}
	fields = append(fields, "event", event)

	if t.err != nil {
		fields = append(fields, "err", t.err)
	}

	if t.fields != nil {
		fields = append(fields, t.fields...)
	}

	return fields
}

func (t *Trace) String() string {
	fields := t.tx.getFields()
	fields = append(fields, t.getFields()...)
	return strings.Join(fields2KVP(fields), " ")
}
