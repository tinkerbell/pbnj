package diagnostic

import (
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	common "github.com/tinkerbell/pbnj/grpc/oob"
)

type Action struct {
	common.Accessory
	ScreenshotRequest          *v1.ScreenshotRequest
	ClearSystemEventLogRequest *v1.ClearSystemEventLogRequest
	SendNMIRequest             *v1.SendNMIRequest
	SystemEventLogRequest      *v1.SystemEventLogRequest
	SystemEventLogRawRequest   *v1.SystemEventLogRawRequest
	ActionName                 string
	RPCName                    string
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

// WithLabels adds the custom tracing and logging labels to an Action struct.
func WithLabels(actionName string, rpcName string) Option {
	return func(a *Action) error {
		a.ActionName = actionName
		a.RPCName = rpcName
		return nil
	}
}

// Option to add to an Actions.
type Option func(a *Action) error
