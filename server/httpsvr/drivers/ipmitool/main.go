// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"github.com/tinkerbell/pbnj/server/httpsvr/evlog"
	"github.com/tinkerbell/pbnj/server/httpsvr/log"
)

var (
	logger log.Logger
	elog   *evlog.Log
)

func SetupLogging(l log.Logger) {
	logger = l.Package("power")
	elog = evlog.New(logger)
}
