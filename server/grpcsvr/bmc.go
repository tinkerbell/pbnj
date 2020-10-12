package grpcsvr

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
)

type bmcService struct {
	log        logging.Logger
	taskRunner task.Task
}

func (b *bmcService) networkSource(ctx context.Context, in *v1.NetworkSourceRequest) (*v1.NetworkSourceResponse, error) {
	l := b.log.GetContextLogger(ctx)
	l.V(0).Info("setting network source")

	switch in.GetAuthn().Authn.(type) {
	case *v1.Authn_ExternalAuthn:
		l.V(1).Info("using external authn")
	default:
		l.V(1).Info("using direct authn")
	}

	return &v1.NetworkSourceResponse{
		TaskId: "good",
	}, nil
}

func (b *bmcService) resetAction(ctx context.Context, in *v1.ResetRequest) (*v1.ResetResponse, error) {
	l := b.log.GetContextLogger(ctx)
	l.V(0).Info("reset action")

	switch in.GetAuthn().Authn.(type) {
	case *v1.Authn_ExternalAuthn:
		l.V(1).Info("using external authn")
	default:
		l.V(1).Info("using direct authn")
	}

	return &v1.ResetResponse{
		TaskId: "good",
	}, nil
}
