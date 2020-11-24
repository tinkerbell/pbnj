package http

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tinkerbell/pbnj/server/grpcsvr/taskrunner"
)

var (
	logger     logr.Logger
	taskRunner *taskrunner.Runner
)

func WithLogger(log logr.Logger) {
	logger = log
}

func WithTaskRunner(runner *taskrunner.Runner) {
	taskRunner = runner
}

func RegisterHandlers() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthcheck", handleHealthcheck)
	http.HandleFunc("/_/ready", handleReady)
	http.HandleFunc("/_/live", handleLive)
}

func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}
