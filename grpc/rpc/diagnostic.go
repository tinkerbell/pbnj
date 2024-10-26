package rpc

import (
	"context"
	"time"

	"github.com/bmc-toolbox/bmclib/v2/bmc"
	"github.com/rs/xid"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/oob/diagnostic"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (d *DiagnosticService) SendNMI(ctx context.Context, in *v1.SendNMIRequest) (*emptypb.Empty, error) {
	empty := &emptypb.Empty{}
	l := logging.ExtractLogr(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)

	l.Info(
		"start Send NMI request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"host", in.Authn.GetDirectAuthn().GetHost(),
	)

	action, err := diagnostic.NewNMISender(in, diagnostic.WithLogger(l))
	if err != nil {
		l.Error(err, "error creating NMI sender")
		return empty, err
	}

	err = action.SendNMI(ctx)
	if err != nil {
		l.Error(err, "error sending NMI")
		return empty, err
	}

	return empty, nil
}

func (d *DiagnosticService) SystemEventLog(ctx context.Context, in *v1.SystemEventLogRequest) (*v1.SystemEventLogResponse, error) {
	l := logging.ExtractLogr(ctx)

	l.Info("start Get System Event Log request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
	)

	selaction, err := diagnostic.NewSystemEventLogAction(in, diagnostic.WithLogger(l),
		diagnostic.WithLabels("system_event_log", "SystemEventLog"))
	if err != nil {
		l.Error(err, "error creating system event log action")
		return nil, err
	}

	entries, _, err := selaction.SystemEventLog(ctx)
	if err != nil {
		l.Error(err, "error getting system event log")
		return nil, err
	}

	events := convertEntriesToEvents(entries)

	return &v1.SystemEventLogResponse{
		Events: events,
	}, nil
}

func (d *DiagnosticService) SystemEventLogRaw(ctx context.Context, in *v1.SystemEventLogRawRequest) (*v1.SystemEventLogRawResponse, error) {
	l := logging.ExtractLogr(ctx)

	l.Info("start Get System Event Log Raw request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
	)

	rawselaction, err := diagnostic.NewSystemEventLogAction(in, diagnostic.WithLogger(l),
		diagnostic.WithLabels("system_event_log_raw", "SystemEventLogRaw"))
	if err != nil {
		l.Error(err, "error creating raw system event log action")
		return nil, err
	}

	_, eventlog, err := rawselaction.SystemEventLog(ctx)
	if err != nil {
		l.Error(err, "error getting raw system event log")
		return nil, err
	}

	return &v1.SystemEventLogRawResponse{
		Log: eventlog,
	}, nil
}

func convertEntriesToEvents(entries bmc.SystemEventLogEntries) []*v1.SystemEventLogEntry {
	var events []*v1.SystemEventLogEntry

	if len(entries) == 0 {
		return events
	}

	for _, entry := range entries {
		if len(entry) < 4 {
			continue
		}
		events = append(events, &v1.SystemEventLogEntry{
			Id:          entry[0],
			Timestamp:   entry[1],
			Description: entry[2],
			Message:     entry[3],
		})
	}

	return events
}
