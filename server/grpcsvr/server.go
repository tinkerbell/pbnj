package grpcsvr

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/healthcheck"
	"github.com/tinkerbell/pbnj/pkg/http"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/repository"
	"github.com/tinkerbell/pbnj/server/grpcsvr/persistence"
	"github.com/tinkerbell/pbnj/server/grpcsvr/rpc"
	"github.com/tinkerbell/pbnj/server/grpcsvr/taskrunner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server options
type Server struct {
	repository.Actions
	bmcTimeout time.Duration
}

// ServerOption for setting optional values
type ServerOption func(*Server)

// WithPersistence sets the log level
func WithPersistence(repo repository.Actions) ServerOption {
	return func(args *Server) { args.Actions = repo }
}

// WithBmcTimeout sets the timeout for BMC calls
func WithBmcTimeout(t time.Duration) ServerOption {
	return func(args *Server) { args.bmcTimeout = t }
}

// RunServer registers all services and runs the server
func RunServer(ctx context.Context, log logging.Logger, grpcServer *grpc.Server, port string, httpServer *http.HTTPServer, opts ...ServerOption) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defaultStore := gokv.Store(freecache.NewStore(freecache.DefaultOptions))

	// instantiate a Repository for task persistence
	repo := &persistence.GoKV{
		Store: defaultStore,
		Ctx:   ctx,
	}

	defaultServer := &Server{
		Actions:    repo,
		bmcTimeout: 15 * time.Second,
	}

	for _, opt := range opts {
		opt(defaultServer)
	}

	taskRunner := &taskrunner.Runner{
		Repository: defaultServer.Actions,
		Ctx:        ctx,
		Log:        log,
	}

	ms := rpc.MachineService{
		Log:        log,
		TaskRunner: taskRunner,
		Timeout:    defaultServer.bmcTimeout,
	}
	v1.RegisterMachineServer(grpcServer, &ms)

	bs := rpc.BmcService{
		Log:        log,
		TaskRunner: taskRunner,
		Timeout:    defaultServer.bmcTimeout,
	}
	v1.RegisterBMCServer(grpcServer, &bs)

	ts := rpc.TaskService{
		Log:        log,
		TaskRunner: taskRunner,
	}
	v1.RegisterTaskServer(grpcServer, &ts)

	grpc_prometheus.Register(grpcServer)

	hc := healthcheck.NewHealthChecker()
	grpc_health_v1.RegisterHealthServer(grpcServer, hc)

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	httpServer.WithTaskRunner(taskRunner)

	go func() {
		err := httpServer.Run()
		if err != nil {
			log.Error(err, "failed to serve http")
			os.Exit(1)
		}
	}()

	// graceful shutdowns
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			log.Info(0, "sig received, shutting down PBnJ")
			grpcServer.GracefulStop()
			<-ctx.Done()
		}
	}()

	go func() {
		<-ctx.Done()
		log.Info(0, "ctx cancelled, shutting down PBnJ")
		grpcServer.GracefulStop()
	}()

	log.Info(0, "starting PBnJ gRPC server")
	return grpcServer.Serve(listen)
}
