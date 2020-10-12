package zaplog

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.SugaredLogger
type Logger struct {
	logr.Logger
	logLevel              string
	outputPaths           []string
	serviceName           string
	keysAndValues         map[string]interface{}
	enableErrLogsToStderr bool
}

// LoggerOption for setting optional values
type LoggerOption func(*Logger)

// WithLogLevel sets the log level
func WithLogLevel(level string) LoggerOption {
	return func(args *Logger) { args.logLevel = level }
}

// WithOutputPaths adds output paths
func WithOutputPaths(paths []string) LoggerOption {
	return func(args *Logger) { args.outputPaths = paths }
}

// WithServiceName adds a service name a logged field
func WithServiceName(name string) LoggerOption {
	return func(args *Logger) { args.serviceName = name }
}

// WithKeysAndValues adds extra key/value fields
func WithKeysAndValues(kvs map[string]interface{}) LoggerOption {
	return func(args *Logger) { args.keysAndValues = kvs }
}

// WithEnableErrLogsToStderr sends .Error logs to stderr
func WithEnableErrLogsToStderr(enable bool) LoggerOption {
	return func(args *Logger) { args.enableErrLogsToStderr = enable }
}

// GetContextLogger get and return a logger from a ctx
func (l Logger) GetContextLogger(ctx context.Context) logr.Logger {
	return zapr.NewLogger(ctxzap.Extract(ctx))
}

// RegisterLogger returns a logr and a zap logger (needed for use in grpc interceptors)
func RegisterLogger(opts ...LoggerOption) (logging.Logger, *zap.Logger, error) {

	// defaults
	const (
		defaultLogLevel    = "info"
		defaultServiceName = "github.com/tinkerbell/pbnj"
	)
	var (
		defaultOutputPaths   = []string{"stdout"}
		defaultKeysAndValues = map[string]interface{}{"service": defaultServiceName}
		zapConfig            = zap.NewProductionConfig()
		zLevel               = zap.InfoLevel
	)

	l := &Logger{
		Logger:        nil,
		logLevel:      defaultLogLevel,
		outputPaths:   defaultOutputPaths,
		serviceName:   defaultServiceName,
		keysAndValues: defaultKeysAndValues,
	}

	for _, opt := range opts {
		opt(l)
	}

	switch l.logLevel {
	case "debug":
		zLevel = zap.DebugLevel
	}

	zapConfig.Level = zap.NewAtomicLevelAt(zLevel)

	zapConfig.OutputPaths = l.outputPaths
	zapConfig.OutputPaths = sliceDedupe(append(zapConfig.OutputPaths, "stdout"))
	zapConfig.InitialFields = l.keysAndValues
	zapLogger, err := zapConfig.Build()
	if err != nil {
		return l, zapLogger, errors.Wrap(err, "failed to build logger config")
	}

	if l.enableErrLogsToStderr {
		errorLogs := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		nonErrorLogs := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})
		console := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)
		encoder := zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, console, nonErrorLogs),
			zapcore.NewCore(encoder, consoleErrors, errorLogs),
		)
		splitLogger := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return core
		})
		zapLogger = zapLogger.WithOptions(splitLogger).Named(l.serviceName)

	}

	l.Logger = zapr.NewLogger(zapLogger)
	return l, zapLogger, nil
}

func sliceDedupe(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}
