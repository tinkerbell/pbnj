package rpc

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
)

// BmcService for doing BMC actions
type BmcService struct {
	Log        logging.Logger
	TaskRunner task.Task
	v1.UnimplementedBMCServer
}

// NetworkSource sets the BMC network source
func (b *BmcService) NetworkSource(ctx context.Context, in *v1.NetworkSourceRequest) (*v1.NetworkSourceResponse, error) {
	l := b.Log.GetContextLogger(ctx)
	l.V(0).Info("setting network source")

	return &v1.NetworkSourceResponse{
		TaskId: "good",
	}, nil
}

// Reset calls a reset on a BMC
func (b *BmcService) Reset(ctx context.Context, in *v1.ResetRequest) (*v1.ResetResponse, error) {
	l := b.Log.GetContextLogger(ctx)
	l.V(0).Info("reset action")

	return &v1.ResetResponse{
		TaskId: "good",
	}, nil
}
