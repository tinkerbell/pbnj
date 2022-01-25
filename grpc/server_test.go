package grpc

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/packethost/pkg/log/logr/v2"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	"github.com/tinkerbell/pbnj/grpc/persistence"
	"github.com/tinkerbell/pbnj/pkg/http"
	"github.com/tinkerbell/pbnj/pkg/zaplog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func TestRunServer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 9*time.Second)
	logger, zapLogger, err := logr.NewPacketLogr()
	log := zaplog.RegisterLogger(logger.Logger)
	ctx = ctxzap.ToContext(ctx, zapLogger)
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	min := 40041
	max := 40042
	port := rand.Intn(max-min+1) + min

	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	defer s.Close()
	repo := &persistence.GoKV{Store: s, Ctx: ctx}

	grpcServer := grpc.NewServer()
	httpServer := http.NewServer(fmt.Sprintf(":%d", port+1))
	httpServer.WithLogger(logger.Logger)

	g := new(errgroup.Group)
	g.Go(func() error {
		return RunServer(ctx, log, grpcServer, strconv.Itoa(port), httpServer, WithPersistence(repo))
	})

	time.Sleep(500 * time.Millisecond)
	cancel()
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestRunServerSignals(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	logger, zapLogger, err := logr.NewPacketLogr()
	log := zaplog.RegisterLogger(logger.Logger)
	ctx = ctxzap.ToContext(ctx, zapLogger)
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	min := 40044
	max := 40045
	port := rand.Intn(max-min+1) + min
	grpcServer := grpc.NewServer()
	httpServer := http.NewServer(fmt.Sprintf(":%d", port+1))
	httpServer.WithLogger(logger.Logger)

	g := new(errgroup.Group)
	g.Go(func() error {
		return RunServer(ctx, log, grpcServer, strconv.Itoa(port), httpServer)
	})

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	_ = proc.Signal(os.Interrupt)

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestRunServerPortInUse(t *testing.T) {
	port := 40039

	// listen on a port
	test, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		t.Fatal(err)
	}
	defer test.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	logger, zapLogger, err := logr.NewPacketLogr()
	log := zaplog.RegisterLogger(logger.Logger)
	ctx = ctxzap.ToContext(ctx, zapLogger)
	if err != nil {
		t.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	httpServer := http.NewServer(fmt.Sprintf(":%d", port+1))
	httpServer.WithLogger(logger.Logger)

	err = RunServer(ctx, log, grpcServer, strconv.Itoa(port), httpServer)
	if err.Error() != "listen tcp :40039: bind: address already in use" {
		t.Fatal(err)
	}
}
