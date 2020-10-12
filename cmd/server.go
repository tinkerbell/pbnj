package cmd

import (
	"context"
	"fmt"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/spf13/cobra"
	"github.com/tinkerbell/pbnj/cmd/zaplog"
	"github.com/tinkerbell/pbnj/server/grpcsvr"

	"goa.design/goa/grpc/middleware"
	"google.golang.org/grpc"
)

const (
	requestIDKey    = "x-request-id"
	requestIDLogKey = "requestID"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run PBnJ server",
	Long:  `Run PBnJ server for interacting with BMCs.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		logger, zlog, err := zaplog.RegisterLogger(
			zaplog.WithServiceName("github.com/tinkerbell/pbnj"),
			zaplog.WithLogLevel(logLevel),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer zlog.Sync() // nolint

		/*
			logger, _, unary, err := packetlog.RegisterPacketLogger("github.com/tinkerbell/pbnj")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		*/

		// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
		grpc_zap.ReplaceGrpcLoggerV2(zlog)

		grpcServer := grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				middleware.UnaryRequestID(middleware.UseXRequestIDMetadataOption(true), middleware.XRequestMetadataLimitOption(512)),
				zaplog.UnaryLogRequestID(zlog, requestIDKey, requestIDLogKey),
				grpc_zap.UnaryServerInterceptor(zlog),
				//unary,
				grpc_validator.UnaryServerInterceptor(),
			),
		)

		// reference for using another persistence layer
		/*
			opts := consul.Options{
				Scheme:  "http",
				Address: "localhost:8500",
				Codec:   encoding.JSON,
			}
			store, err := consul.NewClient(opts)
			if err != nil {
				logger.V(0).Error(err, "failed to connect to consul")
				os.Exit(1)
			}
		*/
		// add grpcsvc.WithPersistence(store) to grpc.RunServer

		if err := grpcsvr.RunServer(ctx, logger, grpcServer, "50051"); err != nil {
			logger.Error(err, "error running server")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
