package http

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tinkerbell/pbnj/grpc/taskrunner"
)

type Server struct {
	address    string
	logger     logr.Logger
	mux        *http.ServeMux
	taskRunner *taskrunner.Runner
}

func (h *Server) WithLogger(log logr.Logger) *Server {
	h.logger = log
	return h
}

func (h *Server) WithTaskRunner(runner *taskrunner.Runner) *Server {
	h.taskRunner = runner
	return h
}

func (h *Server) init() {
	h.mux = http.NewServeMux()
	h.mux.Handle("/metrics", promhttp.Handler())
	h.mux.HandleFunc("/healthcheck", h.handleHealthcheck)
	h.mux.HandleFunc("/_/ready", h.handleReady)
	h.mux.HandleFunc("/_/live", h.handleLive)
}

func (h *Server) Run() error {
	return http.ListenAndServe(h.address, h.mux)
}

func NewServer(addr string) *Server {
	server := &Server{
		address: addr,
	}
	server.init()
	return server
}
