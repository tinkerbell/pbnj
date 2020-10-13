package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/tinkerbell/pbnj/cmd"
	"github.com/tinkerbell/pbnj/cmd/zaplog"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/server/grpcsvr"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"goa.design/goa/grpc/middleware"
	"google.golang.org/grpc"
)

const (
	requestIDKey    = "x-request-id"
	requestIDLogKey = "requestID"
)

type logMessage struct {
	Level   string              `json:"level"`
	Ts      float64             `json:"ts"`
	Caller  string              `json:"caller"`
	Msg     string              `json:"msg"`
	Service string              `json:"service"`
	Resp    []v1.StatusResponse `json:"resp"`
}

func main() {
	port := "40041"

	go func() {
		runGrpcServer(port)
		//serverCmd := cmd.NewRootCmd()
		//serverCmd.SetArgs([]string{"server", "--port", port})
		//serverCmd.Execute()
	}()

	time.Sleep(2 * time.Second)
	clientOutput := captureOutput(func() {
		clientCmd := cmd.NewRootCmd()
		clientCmd.SetArgs([]string{"client", "machine", "--port", port})
		clientCmd.Execute()
	})

	clientStruct := new(logMessage)
	err := json.Unmarshal([]byte(strings.TrimSuffix(clientOutput, "\n")), clientStruct)
	if err != nil {
		fmt.Println("err: " + err.Error())
	}
	fmt.Printf("clientStruct: %+v\n", clientStruct)

	fmt.Println("end")

}

func setupLogsCapture() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	return zap.New(core), logs
}

func runGrpcServer(port string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	zlog, logs := setupLogsCapture()
	l.Logger = zapr.NewLogger(zlog)

	logger, zlog, err := zaplog.RegisterLogger(
		zaplog.WithServiceName("github.com/tinkerbell/pbnj"),
		zaplog.WithOutputPaths([]string{"server.log"}),
	)
	if err != nil {
		return err
	}
	defer zlog.Sync() // nolint

	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLoggerV2(zlog)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			middleware.UnaryRequestID(middleware.UseXRequestIDMetadataOption(true), middleware.XRequestMetadataLimitOption(512)),
			zaplog.UnaryLogRequestID(zlog, requestIDKey, requestIDLogKey),
			grpc_zap.UnaryServerInterceptor(zlog),
			grpc_validator.UnaryServerInterceptor(),
		),
	)

	logger.V(0).Info("serving pbnj ", "port: ", port)
	if err := grpcsvr.RunServer(ctx, logger, grpcServer, port); err != nil {
		return err
	}
	return nil
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = writer
	os.Stderr = writer
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
