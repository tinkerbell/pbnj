package rpc

import (
	"context"
	"time"

	"github.com/rs/xid"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/oob/diagnostic"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
	"go.opentelemetry.io/otel/trace"
)

type DiagnosticService struct {
	v1.UnimplementedDiagnosticServer
	TaskRunner task.Task
	Timeout    time.Duration
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

func (d *DiagnosticService) ClearSystemEventLog(ctx context.Context, in *v1.ClearSystemEventLogRequest) (*v1.ClearSystemEventLogResponse, error) {
	l := logging.ExtractLogr(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)

	l.Info(
		"start Clear System Event Log request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
	)

	execFunc := func(s chan string) (string, error) {
		csl, err := diagnostic.NewSystemEventLogClearer(
			in,
			diagnostic.WithLogger(l),
			diagnostic.WithStatusMessage(s),
		)
		if err != nil {
			return "", err
		}
		// Because this is a background task, we want to pass through the span context, but not be
		// a child context. This allows us to correctly plumb otel into the background task.
		c := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
		taskCtx, cancel := context.WithTimeout(c, d.Timeout)
		defer cancel()
		return csl.ClearSystemEventLog(taskCtx)
	}

	d.TaskRunner.Execute(ctx, l, "clearing system event log", taskID, execFunc)

	return &v1.ClearSystemEventLogResponse{TaskId: taskID}, nil
}
