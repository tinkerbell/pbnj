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

func NewSystemEventLogGetter(req *v1.GetSystemEventLogRequest, opts ...Option) (*Action, error) {
	a := &Action{}
	a.GetSystemEventLogRequest = req
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func NewSystemEventLogRawGetter(req *v1.GetSystemEventLogRawRequest, opts ...Option) (*Action, error) {
	a := &Action{}
	a.GetSystemEventLogRawRequest = req
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (m Action) GetSystemEventLog(ctx context.Context) (result bmc.SystemEventLogEntries, err error) {
	labels := prometheus.Labels{
		"service": "diagnostic",
		"action":  "get_system_event_log",
	}

	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "diagnostic.GetSystemEventLog", trace.WithAttributes(
		attribute.String("bmc.device", m.GetSystemEventLogRequest.GetAuthn().GetDirectAuthn().GetHost().GetHost()),
	))
	defer span.End()

	if v := m.GetSystemEventLogRequest.GetVendor(); v != nil {
		span.SetAttributes(attribute.String("bmc.vendor", v.GetName()))
	}

	host, user, password, parseErr := m.ParseAuth(m.GetSystemEventLogRequest.GetAuthn())
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
	client.Registry.Drivers = client.Registry.Supports(providers.FeatureGetSystemEventLog)

	m.SendStatusMessage("connecting to BMC")
	err = client.Open(ctx)
	meta := client.GetMetadata()
	span.SetAttributes(attribute.StringSlice("bmc.open.providersAttempted", meta.ProvidersAttempted),
		attribute.StringSlice("bmc.open.successfulOpenConns", meta.SuccessfulOpenConns))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, &repository.Error{
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

	// Get the system event log
	m.SendStatusMessage("getting system event log")
	sel, err := client.GetSystemEventLog(ctx)
	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	meta = client.GetMetadata()
	span.SetAttributes(attribute.StringSlice("bmc.get_system_event_log.providersAttempted", meta.ProvidersAttempted),
		attribute.StringSlice("bmc.get_system_event_log.successfulOpenConns", meta.SuccessfulOpenConns))
	if err != nil {
		log.Error(err, "error getting system event log")
		span.SetStatus(codes.Error, "error getting system event log: "+err.Error())
		m.SendStatusMessage(fmt.Sprintf("failed to get system event log %v", host))

		return nil, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	span.SetStatus(codes.Ok, "")
	log.Info("got system event log", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage(fmt.Sprintf("got system event log on %v", host))

	return sel, nil
}

func (m Action) GetSystemEventLogRaw(ctx context.Context) (result string, err error) {
	labels := prometheus.Labels{
		"service": "diagnostic",
		"action":  "get_system_event_log_raw",
	}

	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "diagnostic.GetSystemEventLogRaw", trace.WithAttributes(
		attribute.String("bmc.device", m.GetSystemEventLogRawRequest.GetAuthn().GetDirectAuthn().GetHost().GetHost()),
	))
	defer span.End()

	if v := m.GetSystemEventLogRawRequest.GetVendor(); v != nil {
		span.SetAttributes(attribute.String("bmc.vendor", v.GetName()))
	}

	host, user, password, parseErr := m.ParseAuth(m.GetSystemEventLogRawRequest.GetAuthn())
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
	client.Registry.Drivers = client.Registry.Supports(providers.FeatureGetSystemEventLogRaw)

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

	// Get the system event log
	m.SendStatusMessage("getting system event log")
	sel, err := client.GetSystemEventLogRaw(ctx)
	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	meta = client.GetMetadata()
	span.SetAttributes(attribute.StringSlice("bmc.get_system_event_log_raw.providersAttempted", meta.ProvidersAttempted),
		attribute.StringSlice("bmc.get_system_event_log_raw.successfulOpenConns", meta.SuccessfulOpenConns))
	if err != nil {
		log.Error(err, "error getting raw system event log")
		span.SetStatus(codes.Error, "error getting raw system event log: "+err.Error())
		m.SendStatusMessage(fmt.Sprintf("failed to get raw system event log %v", host))

		return "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	span.SetStatus(codes.Ok, "")
	log.Info("got raw system event log", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage(fmt.Sprintf("got raw system event log on %v", host))

	return sel, nil
}
