package zaplog

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/tinkerbell/pbnj/pkg/logging"
)

// Logger is a wrapper around zap.SugaredLogger.
type Logger struct {
	logr.LogSink
}

// GetContextLogger get and return a logger from a ctx.
func (l Logger) GetContextLogger(ctx context.Context) logr.Logger {
	return zapr.NewLogger(ctxzap.Extract(ctx))
}

// RegisterLogger returns a logr and a zap logger (needed for use in grpc interceptors).
func RegisterLogger(log logr.Logger) logging.Logger {
	return Logger{log.GetSink()}
}
