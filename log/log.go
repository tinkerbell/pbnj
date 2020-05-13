// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

// package log sets up a shared zap.Logger that can be used by all packages.
package log

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	logger   *zap.Logger
	logLevel = zap.LevelFlag("log-level", zap.InfoLevel, "Log level, one of FATAL, PANIC, DPANIC, ERROR, WARN, INFO, or DEBUG")
)

// Logger is a wrapper around zap.SugaredLogger
type Logger struct {
	*zap.SugaredLogger
}

// Init initializes the logging system and sets the "service" key to the provided argument.
// This func should only be called once and after flag.Parse() has been called otherwise leveled logging will not be configured correctly.
func Init(service string) (Logger, func() error, error) {
	var config zap.Config
	if os.Getenv("DEBUG") != "" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.Level = zap.NewAtomicLevelAt(*logLevel)

	l, err := config.Build()
	if err != nil {
		err = errors.Wrap(err, "failed to build logger config")
	}
	logger = l.With(zap.String("service", service))

	if err != nil {
		return Logger{}, nil, err
	}
	return Logger{logger.Sugar()}, logger.Sync, nil
}

func (l Logger) Notice(args ...interface{}) {
	l.AddCallerSkip(1).Info(args...)
}

func (l Logger) Trace(args ...interface{}) {
	l.AddCallerSkip(1).Debug(args...)
}

func (l Logger) Warning(args ...interface{}) {
	l.AddCallerSkip(1).Warn(args...)
}

func (l Logger) With(args ...interface{}) Logger {
	return Logger{l.SugaredLogger.With(args...)}
}

func (l Logger) AddCallerSkip(skip int) Logger {
	s := l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
	return Logger{s}
}

func (l Logger) Package(pkg string) Logger {
	return Logger{l.SugaredLogger.With("pkg", pkg)}
}
