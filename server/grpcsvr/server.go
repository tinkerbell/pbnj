package grpcsvr

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/repository"
	"github.com/tinkerbell/pbnj/pkg/task"
	"github.com/tinkerbell/pbnj/server/grpcsvr/persistence"
	"github.com/tinkerbell/pbnj/server/grpcsvr/taskrunner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server options
type Server struct {
	Persistence gokv.Store
}

// ServerOption for setting optional values
type ServerOption func(*Server)

// WithPersistence sets the log level
func WithPersistence(store gokv.Store) ServerOption {
	return func(args *Server) { args.Persistence = store }
}

// RunServer registers all services and runs the server
func RunServer(ctx context.Context, log logging.Logger, grpcServer *grpc.Server, port string, opts ...ServerOption) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defaultServer := &Server{
		Persistence: gokv.Store(freecache.NewStore(freecache.DefaultOptions)),
	}
	for _, opt := range opts {
		opt(defaultServer)
	}

	// instantiate a Repository for task persistence
	var repo repository.Actions
	repo = &persistence.GoKV{
		Store: defaultServer.Persistence,
		Ctx:   ctx,
	}

	var taskRunner task.Task
	taskRunner = &taskrunner.Runner{
		Repository: repo,
		Ctx:        ctx,
		Log:        log,
	}

	ms := machineService{
		log:        log,
		taskRunner: taskRunner,
	}
	m := v1.MachineService{
		BootDevice: ms.device,
		Power:      ms.powerAction,
	}
	v1.RegisterMachineService(grpcServer, &m)

	bs := bmcService{
		log:        log,
		taskRunner: taskRunner,
	}
	b := v1.BMCService{
		NetworkSource: bs.networkSource,
		Reset:         bs.resetAction,
	}
	v1.RegisterBMCService(grpcServer, &b)

	ts := taskService{
		log:        log,
		taskRunner: taskRunner,
	}
	t := v1.TaskService{
		Status: ts.Task,
	}
	v1.RegisterTaskService(grpcServer, &t)

	us := userService{
		log:        log,
		taskRunner: taskRunner,
	}
	u := v1.UserService{
		CreateUser: us.createUser,
		DeleteUser: us.deleteUser,
		UpdateUser: us.updateUser,
	}
	v1.RegisterUserService(grpcServer, &u)
	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// graceful shutdowns
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			log.V(0).Info("sig received, shutting down PBnJ")
			grpcServer.GracefulStop()
			defaultServer.Persistence.Close()
			<-ctx.Done()
		}
	}()

	go func() {
		<-ctx.Done()
		log.V(0).Info("ctx cancelled, shutting down PBnJ")
		grpcServer.GracefulStop()
		defaultServer.Persistence.Close()
	}()

	log.V(0).Info("starting PBnJ gRPC server")
	return grpcServer.Serve(listen)
}
