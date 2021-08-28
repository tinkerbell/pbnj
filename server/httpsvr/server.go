// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package httpsvr

import (
	"flag"
	"os"
	"time"

	"github.com/tinkerbell/pbnj/server/httpsvr/api"
	"github.com/tinkerbell/pbnj/server/httpsvr/drivers/ipmitool"
	"github.com/tinkerbell/pbnj/server/httpsvr/drivers/racadm"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/power"
	"github.com/tinkerbell/pbnj/server/httpsvr/log"
	"github.com/tinkerbell/pbnj/server/httpsvr/reqid"
)

var (
	GitRev     = "unknown (use make)"
	listenAddr = ":9090"
)

func init() {
	if nomadPort, ok := os.LookupEnv("NOMAD_PORT_internal_http"); ok {
		listenAddr = ":" + nomadPort
	}

	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "IP and port to listen on for HTTP")
}

func cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		power.CleanupTasks(3 * time.Minute)
	}
}

// RunHTTPServer runs PBnJ v1 HTTP server.
func RunHTTPServer() {
	flag.Parse()

	logger, sync, err := log.Init("github.com/tinkerbell/pbnj")
	if err != nil {
		panic(err)
	}

	logger = logger.Package("main")
	defer func() {
		err := sync()
		if err != nil {
			logger.Errorf("logger sync failed: %v", err)
		}
	}()

	api.SetupLogging(logger)
	ipmitool.SetupLogging(logger)
	power.SetupLogging(logger)
	racadm.SetupLogging(logger)
	reqid.SetupLogging(logger)
	go cleanup()

	if err := api.Serve(listenAddr, GitRev); err != nil {
		logger.Fatalw("error serving api", "error", err)
	}
}
