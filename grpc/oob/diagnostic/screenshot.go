package diagnostic

import (
	"context"
	"fmt"
	"time"

	"github.com/bmc-toolbox/bmclib/v2"
	"github.com/bmc-toolbox/bmclib/v2/bmc"
	"github.com/bmc-toolbox/bmclib/v2/providers"
	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	common "github.com/tinkerbell/pbnj/grpc/oob"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func NewScreenshotter(req *v1.ScreenshotRequest, opts ...Option) (*Action, error) {
	a := &Action{}
	a.ScreenshotRequest = req
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (m Action) GetScreenshot(ctx context.Context) (image []byte, filetype string, err error) {
	labels := prometheus.Labels{
		"service": "diagnostic",
		"action":  "screenshot",
	}

	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "diagnostic.GetScreenshot", trace.WithAttributes(
		attribute.String("bmc.device", m.ScreenshotRequest.GetAuthn().GetDirectAuthn().GetHost().GetHost()),
	))
	defer span.End()

	if v := m.ScreenshotRequest.GetVendor(); v != nil {
		span.SetAttributes(attribute.String("bmc.vendor", v.GetName()))
	}

	host, user, password, parseErr := m.ParseAuth(m.ScreenshotRequest.GetAuthn())
	if parseErr != nil {
		span.SetStatus(codes.Error, "error parsing credentials: "+parseErr.Error())
		return nil, "", parseErr
	}
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	opts := []bmclib.Option{
		bmclib.WithLogger(m.Log),
		bmclib.WithPerProviderTimeout(common.BMCTimeoutFromCtx(ctx)),
	}

	client := bmclib.NewClient(host, user, password, opts...)
	client.Registry.Drivers = client.Registry.Supports(providers.FeatureScreenshot)

	m.SendStatusMessage("connecting to BMC")
	err = client.Open(ctx)
	meta := client.GetMetadata()
	span.SetAttributes(attribute.StringSlice("bmc.open.providersAttempted", meta.ProvidersAttempted),
		attribute.StringSlice("bmc.open.successfulOpenConns", meta.SuccessfulOpenConns))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, "", &repository.Error{
			Code:    v1.Code_value["PERMISSION_DENIED"],
			Message: err.Error(),
		}
	}
	log := m.Log.WithValues("host", host, "user", user)
	defer func() {
		client.Close(ctx)
		log.Info("closed connections", logMetadata(client.GetMetadata())...)
	}()
	log.Info("connected to BMC", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage("connected to BMC")

	image, filetype, err = client.Screenshot(ctx)
	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	meta = client.GetMetadata()
	span.SetAttributes(attribute.String("bmc.screenshot.successfulProvider", meta.SuccessfulProvider),
		attribute.StringSlice("bmc.screenshot.ProvidersAttempted", meta.ProvidersAttempted))
	if err != nil {
		log.Error(err, "error getting screenshot")
		span.SetStatus(codes.Error, "error getting screenshot: "+err.Error())
		m.SendStatusMessage(fmt.Sprintf("failed to screenshot %v", host))

		return nil, "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	span.SetStatus(codes.Ok, "")
	log.Info("got screenshot", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage(fmt.Sprintf("got screenshot from %v", host))

	return image, filetype, nil
}

func logMetadata(md bmc.Metadata) []interface{} {
	kvs := []interface{}{
		"ProvidersAttempted", md.ProvidersAttempted,
		"SuccessfulOpenConns", md.SuccessfulOpenConns,
		"SuccessfulCloseConns", md.SuccessfulCloseConns,
		"SuccessfulProvider", md.SuccessfulProvider,
	}

	return kvs
}
