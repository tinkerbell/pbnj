package diagnostic

import (
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	common "github.com/tinkerbell/pbnj/grpc/oob"
)

type Action struct {
	common.Accessory
	ScreenshotRequest *v1.ScreenshotRequest
	ClearSELRequest   *v1.ClearSELRequest
}

// WithLogger adds a logr to an Action struct.
func WithLogger(l logr.Logger) Option {
	return func(a *Action) error {
		a.Log = l
		return nil
	}
}

// WithStatusMessage adds a status message chan to an Action struct.
func WithStatusMessage(s chan string) Option {
	return func(a *Action) error {
		a.StatusMessages = s
		return nil
	}
}

// Option to add to an Actions.
type Option func(a *Action) error
