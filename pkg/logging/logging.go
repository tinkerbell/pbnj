package logging

import (
	"context"

	"github.com/go-logr/logr"
)

// Logger represent common interface for logging function
type Logger interface {
	logr.Logger
	GetContextLogger(ctx context.Context) logr.Logger
}
