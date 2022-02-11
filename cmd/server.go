package cmd

import (
	"context"
	"errors"
	"os"
	"time"

	jwt "github.com/cristalhq/jwt/v3"
	"github.com/equinix-labs/otel-init-go/otelinit"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	jwt_helper "github.com/golang-jwt/jwt/v4"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/packethost/pkg/grpc/authz"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	grpcsvr "github.com/tinkerbell/pbnj/grpc"
	"github.com/tinkerbell/pbnj/pkg/http"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"goa.design/goa/grpc/middleware"
	"google.golang.org/grpc"
)

const (
	requestIDKey    = "x-request-id"
	requestIDLogKey = "requestID"
)

var (
	port        string
	metricsAddr string
	enableAuthz bool
	hsKey       string
	rsPubKey    string
	// bmcTimeout is how long a BMC call/interaction is allow to run before it is cancelled.
	bmcTimeout time.Duration
	// serverCmd represents the server command.
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run PBnJ server",
		Long:  `Run PBnJ server for interacting with BMCs.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			ctx, otelShutdown := otelinit.InitOpenTelemetry(ctx, "pbnj")
			defer otelShutdown(ctx)

			logger := defaultLogger(logLevel).WithName("github.com/tinkerbell/pbnj").WithValues("service", "github.com/tinkerbell/pbnj")

			authzInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				return handler(ctx, req)
			}
			if enableAuthz {
				if hsKey != "" || rsPubKey != "" {
					authzInterceptor = grpc_auth.UnaryServerInterceptor(authFunc())
				} else {
					logger.Error(errors.New("error configuring server"), "authorization enabled but no symmetric or asymmetric key was provided")
					os.Exit(1)
				}
			}
			grpc_prometheus.EnableHandlingTimeHistogram()
			grpcServer := grpc.NewServer(
				grpc_middleware.WithUnaryServerChain(
					grpc_prometheus.UnaryServerInterceptor,
					authzInterceptor,
					middleware.UnaryRequestID(middleware.UseXRequestIDMetadataOption(true), middleware.XRequestMetadataLimitOption(512)),
					logging.UnaryServerInterceptor(logger),                   // this puts the logger in the context. Allows per-request logging and other middleware to be used.
					logging.UnaryLogRequestID(requestIDKey, requestIDLogKey), // must be after logging.UnaryServerInterceptor because the logger must be in the context.
					logging.UnaryLogBMCIP(),                                  // must be after logging.UnaryServerInterceptor because the logger must be in the context.
					otelgrpc.UnaryServerInterceptor(),
					grpc_validator.UnaryServerInterceptor(),
				),
			)

			httpServer := http.NewServer(metricsAddr)
			httpServer.WithLogger(logger)

			if err := grpcsvr.RunServer(ctx, logger, grpcServer, port, httpServer, grpcsvr.WithBmcTimeout(bmcTimeout)); err != nil {
				logger.Error(err, "error running server")
				os.Exit(1)
			}
		},
	}
)

func init() {
	serverCmd.PersistentFlags().StringVar(&port, "port", "50051", "grpc server port")
	serverCmd.PersistentFlags().StringVar(&metricsAddr, "metricsListenAddr", ":8080", "metrics server listen address")
	serverCmd.PersistentFlags().BoolVar(&enableAuthz, "enableAuthz", false, "enable Authz middleware. Configure with configuration file details")
	serverCmd.PersistentFlags().StringVar(&hsKey, "hsKey", "", "HS key")
	serverCmd.PersistentFlags().StringVar(&rsPubKey, "rsPubKey", "", "RS public key")
	serverCmd.PersistentFlags().DurationVar(&bmcTimeout, "bmcTimeout", (15 * time.Second), "Timeout for BMC calls")
	rootCmd.AddCommand(serverCmd)
}

// defaultLogger is a zerolog logr implementation.
func defaultLogger(level string) logr.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"

	zl := zerolog.New(os.Stdout)
	zl = zl.With().Caller().Timestamp().Logger()
	var l zerolog.Level
	switch level {
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}
	zl = zl.Level(l)

	return zerologr.New(&zl)
}

// authFunc will validate (signed and not expired) the JWT against the methods in the ScopeMapping.
// No scopes will be checked because scopes can be arbitrary json structures and are generally
// catered to the Authn signing the token. Accepting arbitrary json and using that to validate
// could be a future feature to add if requested.
func authFunc() func(ctx context.Context) (context.Context, error) {
	opts := []authz.ConfigOption{authz.WithDisableAudienceValidation(true)}
	var algo jwt.Algorithm

	if hsKey == "" && rsPubKey == "" {
		return func(ctx context.Context) (context.Context, error) {
			return ctx, errors.New("authorization enabled but no symmetric or asymmetric key was provided")
		}
	}

	if hsKey != "" {
		if rsPubKey != "" {
			return func(ctx context.Context) (context.Context, error) {
				return ctx, errors.New("both an HS Key and an RS public key were provided: choose one")
			}
		}

		algo = jwt.HS256
		opts = append(opts, authz.WithHSKey([]byte(hsKey)))
	}

	if rsPubKey != "" {
		algo = jwt.RS256
		pubKey, err := jwt_helper.ParseRSAPublicKeyFromPEM([]byte(rsPubKey))
		if err != nil {
			return func(ctx context.Context) (context.Context, error) { return ctx, err }
		}
		opts = append(opts, authz.WithRSAPubKey(pubKey))
	}

	protectedMethods := map[string][]string{
		"/github.com.tinkerbell.pbnj.api.v1.Machine/Power":      {},
		"/github.com.tinkerbell.pbnj.api.v1.Machine/BootDevice": {},
		"/github.com.tinkerbell.pbnj.api.v1.BMC/NetworkSource":  {},
		"/github.com.tinkerbell.pbnj.api.v1.BMC/Reset":          {},
		"/github.com.tinkerbell.pbnj.api.v1.BMC/CreateUser":     {},
		"/github.com.tinkerbell.pbnj.api.v1.BMC/DeleteUser":     {},
		"/github.com.tinkerbell.pbnj.api.v1.BMC/UpdateUser":     {},
	}
	config := authz.NewConfig(algo, protectedMethods, opts...)
	return config.AuthFunc
}
