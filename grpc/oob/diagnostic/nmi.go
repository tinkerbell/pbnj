package diagnostic

import (
	"context"
	"fmt"

	"github.com/bmc-toolbox/bmclib/v2"
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

func NewNMISender(req *v1.SendNMIRequest, opts ...Option) (*Action, error) {
	a := &Action{}
	a.SendNMIRequest = req
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (m Action) SendNMI(ctx context.Context) error {
	labels := prometheus.Labels{
		"service": "diagnostic",
		"action":  "send_nmi",
	}

	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "diagnostic.SendNMI", trace.WithAttributes(
		attribute.String("bmc.device", m.SendNMIRequest.GetAuthn().GetDirectAuthn().GetHost().GetHost()),
	))
	defer span.End()

	host, user, password, parseErr := m.ParseAuth(m.SendNMIRequest.GetAuthn())
	if parseErr != nil {
		span.SetStatus(codes.Error, "error parsing credentials: "+parseErr.Error())
		return parseErr
	}
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))

	opts := []bmclib.Option{
		bmclib.WithLogger(m.Log),
		bmclib.WithPerProviderTimeout(common.BMCTimeoutFromCtx(ctx)),
	}

	client := bmclib.NewClient(host, user, password, opts...)

	m.SendStatusMessage("connecting to BMC")
	err := client.Open(ctx)
	meta := client.GetMetadata()
	span.SetAttributes(attribute.StringSlice("bmc.open.providersAttempted", meta.ProvidersAttempted),
		attribute.StringSlice("bmc.open.successfulOpenConns", meta.SuccessfulOpenConns))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return &repository.Error{
			Code:    v1.Code_value["PERMISSION_DENIED"],
			Message: err.Error(),
		}
	}
	log := m.Log.WithValues("host", host, "user", user)
	defer func() {
		_ = client.Close(ctx)
		log.Info("closed connections", logMetadata(client.GetMetadata())...)
	}()
	log.Info("connected to BMC", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage("connected to BMC")

	err = client.SendNMI(ctx)
	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	meta = client.GetMetadata()
	span.SetAttributes(attribute.String("bmc.send_nmi.successfulProvider", meta.SuccessfulProvider),
		attribute.StringSlice("bmc.send_nmi.ProvidersAttempted", meta.ProvidersAttempted))
	if err != nil {
		log.Error(err, "error sending NMI")
		span.SetStatus(codes.Error, "error sending NMI: "+err.Error())
		m.SendStatusMessage(fmt.Sprintf("failed to send NMI %v", host))

		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	span.SetStatus(codes.Ok, "")
	log.Info("NMI sent", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage(fmt.Sprintf("Send NMI to host %v", host))

	return nil
}
