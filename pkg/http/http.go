package http

import (
	"context"
	"net/http"
	"time"

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

func (h *Server) Run(ctx context.Context) error {
	svr := &http.Server{
		Addr:    h.address,
		Handler: h.mux,
		// Mitigate Slowloris attacks. 20 seconds is based on Apache's recommended 20-40
		// recommendation. Hegel doesn't really have many headers so 20s should be plenty of time.
		// https://en.wikipedia.org/wiki/Slowloris_(computer_security)
		ReadHeaderTimeout: 20 * time.Second,
	}

	go func() {
		<-ctx.Done()
		_ = svr.Shutdown(ctx)
	}()

	if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func NewServer(addr string) *Server {
	server := &Server{
		address: addr,
	}
	server.init()
	return server
}
