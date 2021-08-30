// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package evlog

import (
	"fmt"
	"time"
)

// Tx is used for loggin an event.
type Tx struct {
	l      *Log
	ID     string
	fields []interface{}
}

func fields2KVP(fields []interface{}) []string {
	kvp := make([]string, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		pair := fmt.Sprintf("%s=%s", fields[i], fields[i+1])
		kvp = append(kvp, pair)
	}
	return kvp
}

// With returns.
func (tx *Tx) With(fields ...interface{}) *Tx {
	n := *tx
	n.fields = fields
	return &n
}

// Panic is a helper function to log a CRITICAL event.
func (tx *Tx) Panic(event string, fields ...interface{}) {
	tx.l.logger.With(fields...).With(tx.getFields()...).With("event", event).Panic()
}

// Error is a helper function to log an error.
func (tx *Tx) Error(event string, fields ...interface{}) {
	tx.l.logger.With(fields...).With(tx.getFields()...).With("event", event).Error()
}

// Warning is a helper function to log a warning.
func (tx *Tx) Warning(event string, fields ...interface{}) {
	tx.l.logger.With(fields...).With(tx.getFields()...).With("event", event).Warning()
}

// Notice is a helper function to log a notice.
func (tx *Tx) Notice(event string, fields ...interface{}) {
	tx.l.logger.With(fields...).With(tx.getFields()...).With("event", event).Notice()
}

// Info is a helper function to log an info message.
func (tx *Tx) Info(event string, fields ...interface{}) {
	tx.l.logger.With(fields...).With(tx.getFields()...).With("event", event).Info()
}

// Debug is a helper function to log a debug message.
func (tx *Tx) Debug(event string, fields ...interface{}) {
	tx.l.logger.With(fields...).With(tx.getFields()...).With("event", event).Debug()
}

// Trace is a helper function to log a trace message.
func (tx *Tx) Trace(event string, fields ...interface{}) *Trace {
	t := &Trace{
		start:  time.Now(),
		event:  event,
		fields: fields,
		tx:     *tx,
	}
	tx.l.logger.With(fields...).With(t.getFields()...).With(tx.getFields()...).Debug()
	return t
}

func (tx *Tx) getFields() []interface{} {
	if tx.ID == "" {
		return nil
	}
	fields := []interface{}{"txid", tx.ID}

	if tx.fields != nil {
		fields = append(fields, tx.fields...)
	}

	return fields
}
