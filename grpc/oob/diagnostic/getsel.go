package diagnostic

import (
	"context"
	"fmt"

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

func NewSystemEventLogAction(req interface{}, opts ...Option) (*Action, error) {
	a := &Action{}
	switch r := req.(type) {
	case *v1.SystemEventLogRequest:
		a.SystemEventLogRequest = r
	case *v1.SystemEventLogRawRequest:
		a.SystemEventLogRawRequest = r
	default:
		return nil, fmt.Errorf("unsupported request type")
	}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (m Action) SystemEventLog(ctx context.Context) (entries bmc.SystemEventLogEntries, raw string, err error) {
	labels := prometheus.Labels{
		"service": "diagnostic",
		"action":  m.ActionName,
	}

	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "diagnostic."+m.RPCName, trace.WithAttributes(
		attribute.String("bmc.device", m.SystemEventLogRequest.GetAuthn().GetDirectAuthn().GetHost().GetHost()),
	))
	defer span.End()

	if v := m.SystemEventLogRequest.GetVendor(); v != nil {
		span.SetAttributes(attribute.String("bmc.vendor", v.GetName()))
	}

	host, user, password, parseErr := m.ParseAuth(m.SystemEventLogRequest.GetAuthn())
	if parseErr != nil {
		span.SetStatus(codes.Error, "error parsing credentials: "+parseErr.Error())
		return entries, raw, parseErr
	}

	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))

	opts := []bmclib.Option{
		bmclib.WithLogger(m.Log),
		bmclib.WithPerProviderTimeout(common.BMCTimeoutFromCtx(ctx)),
	}

	client := bmclib.NewClient(host, user, password, opts...)

	// Set the driver(s) to use based on the request type
	switch {
	case m.SystemEventLogRequest != nil:
		client.Registry.Drivers = client.Registry.Supports(providers.FeatureGetSystemEventLog)
	case m.SystemEventLogRawRequest != nil:
		client.Registry.Drivers = client.Registry.Supports(providers.FeatureGetSystemEventLogRaw)
	default:
		return entries, raw, fmt.Errorf("unsupported request type")
	}

	m.SendStatusMessage("connecting to BMC")
	err = client.Open(ctx)
	meta := client.GetMetadata()
	span.SetAttributes(attribute.StringSlice("bmc.open.providersAttempted", meta.ProvidersAttempted),
		attribute.StringSlice("bmc.open.successfulOpenConns", meta.SuccessfulOpenConns))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return entries, raw, &repository.Error{
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

	m.SendStatusMessage("getting " + m.ActionName + " on " + host)

	switch {
	case m.SystemEventLogRequest != nil:
		// Get the system event log
		entries, err = client.GetSystemEventLog(ctx)
	case m.SystemEventLogRawRequest != nil:
		// Get the system event log
		raw, err = client.GetSystemEventLogRaw(ctx)
	default:
		return entries, raw, fmt.Errorf("unsupported request type")
	}

	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	meta = client.GetMetadata()
	span.SetAttributes(attribute.StringSlice("bmc."+m.ActionName+".providersAttempted", meta.ProvidersAttempted),
		attribute.StringSlice("bmc."+m.ActionName+".successfulOpenConns", meta.SuccessfulOpenConns))
	if err != nil {
		log.Error(err, "error getting "+m.ActionName)
		span.SetStatus(codes.Error, "error getting "+m.ActionName+": "+err.Error())
		m.SendStatusMessage(fmt.Sprintf("failed to get "+m.ActionName+" %v", host))

		return entries, raw, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	span.SetStatus(codes.Ok, "")
	log.Info("got "+m.ActionName, logMetadata(client.GetMetadata())...)
	m.SendStatusMessage(fmt.Sprintf("got "+m.ActionName+" on %v", host))

	return entries, raw, nil
}
