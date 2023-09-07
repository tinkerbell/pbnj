package diagnostic

import (
	"context"
	"fmt"

	"github.com/bmc-toolbox/bmclib/v2"
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

func NewSELClearer(req *v1.ClearSELRequest, opts ...Option) (*Action, error) {
	a := &Action{}
	a.ClearSELRequest = req
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (m Action) ClearSEL(ctx context.Context) (result string, err error) {
	labels := prometheus.Labels{
		"service": "diagnostic",
		"action":  "clearsel",
	}

	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "diagnostic.ClearSEL", trace.WithAttributes(
		attribute.String("bmc.device", m.ClearSELRequest.GetAuthn().GetDirectAuthn().GetHost().GetHost()),
	))
	defer span.End()

	if v := m.ClearSELRequest.GetVendor(); v != nil {
		span.SetAttributes(attribute.String("bmc.vendor", v.GetName()))
	}

	host, user, password, parseErr := m.ParseAuth(m.ClearSELRequest.GetAuthn())
	if parseErr != nil {
		span.SetStatus(codes.Error, "error parsing credentials: "+parseErr.Error())
		return result, parseErr
	}
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))

	opts := []bmclib.Option{
		bmclib.WithLogger(m.Log),
		bmclib.WithPerProviderTimeout(common.BMCTimeoutFromCtx(ctx)),
	}

	client := bmclib.NewClient(host, user, password, opts...)
	client.Registry.Drivers = client.Registry.Supports(providers.FeatureSELClear)

	m.SendStatusMessage("connecting to BMC")
	err = client.Open(ctx)
	meta := client.GetMetadata()
	span.SetAttributes(attribute.StringSlice("bmc.open.providersAttempted", meta.ProvidersAttempted),
		attribute.StringSlice("bmc.open.successfulOpenConns", meta.SuccessfulOpenConns))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return "", &repository.Error{
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

	err = client.ClearSEL(ctx)
	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	meta = client.GetMetadata()
	span.SetAttributes(attribute.String("bmc.clearsel.successfulProvider", meta.SuccessfulProvider),
		attribute.StringSlice("bmc.clearsel.ProvidersAttempted", meta.ProvidersAttempted))
	if err != nil {
		log.Error(err, "error clearing SEL")
		span.SetStatus(codes.Error, "error clearing SEL: "+err.Error())
		m.SendStatusMessage(fmt.Sprintf("failed to clear SEL %v", host))

		return "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	span.SetStatus(codes.Ok, "")
	log.Info("cleared SEL", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage(fmt.Sprintf("cleared SEL on %v", host))

	return result, nil
}
