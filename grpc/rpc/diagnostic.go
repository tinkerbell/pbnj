package rpc

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/oob/diagnostic"
	"github.com/tinkerbell/pbnj/pkg/logging"
)

type DiagnosticService struct {
	v1.UnimplementedDiagnosticServer
}

func (d *DiagnosticService) Screenshot(ctx context.Context, in *v1.ScreenshotRequest) (*v1.ScreenshotResponse, error) {
	l := logging.ExtractLogr(ctx)

	l = l.WithValues("bmcIP", in.Authn.GetDirectAuthn().GetHost().GetHost())

	l.Info(
		"start Screenshot request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
	)

	ms, err := diagnostic.NewScreenshotter(in, diagnostic.WithLogger(l))
	if err != nil {
		l.Error(err, "error creating screenshotter")
		return nil, err
	}

	image, filetype, err := ms.GetScreenshot(ctx)
	if err != nil {
		l.Error(err, "error getting screenshot")
		return nil, err
	}

	return &v1.ScreenshotResponse{
		Image:    image,
		Filetype: filetype,
	}, nil
}
