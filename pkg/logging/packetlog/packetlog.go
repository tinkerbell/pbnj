// Package packetlog defines an implementation of the github.com/go-logr/logr
// interfaces built on top of packet logging (github.com/packethost/pkg/log).
package packetlog

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/packethost/pkg/log"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"google.golang.org/grpc"
)

// Logger for implementing
type Logger struct {
	logr.Logger
}

// GetContextLogger get and return a logger from a ctx
func (l Logger) GetContextLogger(ctx context.Context) logr.Logger {
	return zapr.NewLogger(ctxzap.Extract(ctx))
}

// packetLogger is a logr.Logger that uses packetLog to log.
type packetLogger struct {
	l *log.Logger
}

func (pl *packetLogger) Enabled() bool {
	return true
}

func (pl *packetLogger) Info(msg string, keysAndVals ...interface{}) {
	nkvs := []interface{}{msg, ""}
	pl.l.Info(append(nkvs, keysAndVals...)...)
}

func (pl *packetLogger) Error(err error, msg string, keysAndVals ...interface{}) {
	nkvs := []interface{}{msg, ""}
	pl.l.Error(err, append(nkvs, keysAndVals...)...)
}

func (pl *packetLogger) V(level int) logr.Logger {
	return &packetLogger{l: pl.l}
}

func (pl *packetLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	newLogger := pl.l.With(keysAndValues...)
	return newLoggerWithExtraSkip(&newLogger, 0)
}

func (pl *packetLogger) WithName(name string) logr.Logger {
	newLogger := pl.l.Package(name)
	return newLoggerWithExtraSkip(&newLogger, 0)
}

func (pl *packetLogger) GRPCLoggers() (grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor) {
	return pl.GRPCLoggers()
}

// newLoggerWithExtraSkip allows creation of loggers with variable levels of callstack skipping
func newLoggerWithExtraSkip(l *log.Logger, callerSkip int) logr.Logger {
	log := l.AddCallerSkip(callerSkip)
	return &packetLogger{l: &log}
}

// NewLogger creates a new logr.Logger using the given packetlogger.
func NewLogger(service string) (logr.Logger, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	newLogger, err := log.Init(service)
	if err != nil {
		return nil, nil, nil, err
	}
	//newLogger := l.AddCallerSkip(0)
	stream, unary := newLogger.GRPCLoggers()
	return &packetLogger{l: &newLogger}, stream, unary, nil
}

// RegisterPacketLogger creates a new logging.Logger
func RegisterPacketLogger(service string) (logging.Logger, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	// creates a new logger skipping one level of callstack
	l, stream, unary, err := NewLogger(service)
	if err != nil {
		return nil, nil, nil, err
	}

	nl := Logger{l}
	return nl, stream, unary, nil
}
