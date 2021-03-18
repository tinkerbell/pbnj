package http

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tinkerbell/pbnj/server/grpcsvr/taskrunner"
)

type HTTPServer struct {
	address    string
	logger     logr.Logger
	mux        *http.ServeMux
	taskRunner *taskrunner.Runner
}

func (h *HTTPServer) WithLogger(log logr.Logger) *HTTPServer {
	h.logger = log
	return h
}

func (h *HTTPServer) WithTaskRunner(runner *taskrunner.Runner) *HTTPServer {
	h.taskRunner = runner
	return h
}

func (h *HTTPServer) init() {
	h.mux = http.NewServeMux()
	h.mux.Handle("/metrics", promhttp.Handler())
	h.mux.HandleFunc("/healthcheck", h.handleHealthcheck)
	h.mux.HandleFunc("/_/ready", h.handleReady)
	h.mux.HandleFunc("/_/live", h.handleLive)
}

func (h *HTTPServer) Run() error {
	return http.ListenAndServe(h.address, h.mux)
}

func NewHTTPServer(addr string) *HTTPServer {
	server := &HTTPServer{
		address: addr,
	}
	server.init()
	return server
}
