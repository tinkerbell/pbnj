// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package evlog

import (
	"context"

	"github.com/tinkerbell/pbnj/log"
	"github.com/tinkerbell/pbnj/reqid"
)

// Log embeds a Logger
type Log struct {
	logger log.Logger
}

// New instantiates a new Log
func New(logger log.Logger) *Log {
	return &Log{logger: logger.AddCallerSkip(1)}
}

// TxFromContext returns a new Tx, using the id from the provided context
func (l *Log) TxFromContext(ctx context.Context) *Tx {
	return &Tx{l: l, ID: reqid.FromContext(ctx)}
}
