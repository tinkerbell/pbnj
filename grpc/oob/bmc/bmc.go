package bmc

import (
	"context"
	"fmt"

	"github.com/bmc-toolbox/bmclib/v2"
	"github.com/bmc-toolbox/bmclib/v2/bmc"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	common "github.com/tinkerbell/pbnj/grpc/oob"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Action for making bmc actions on BMCs, implements oob.User interface.
type Action struct {
	common.Accessory
	CreateUserRequest    *v1.CreateUserRequest
	DeleteUserRequest    *v1.DeleteUserRequest
	UpdateUserRequest    *v1.UpdateUserRequest
	ResetBMCRequest      *v1.ResetRequest
	DeactivateSOLRequest *v1.DeactivateSOLRequest
}

// Option to add to an Actions.
type Option func(a *Action) error

// WithLogger adds a logr to an Action struct.
func WithLogger(l logr.Logger) Option {
	return func(a *Action) error {
		a.Log = l
		return nil
	}
}

// WithStatusMessage adds a status message chan to an Action struct.
func WithStatusMessage(s chan string) Option {
	return func(a *Action) error {
		a.StatusMessages = s
		return nil
	}
}

// WithCreateUserRequest adds CreateUserRequest to an Action struct.
func WithCreateUserRequest(in *v1.CreateUserRequest) Option {
	return func(a *Action) error {
		a.CreateUserRequest = in
		return nil
	}
}

// WithDeleteUserRequest adds DeleteUserRequest to an Action struct.
func WithDeleteUserRequest(in *v1.DeleteUserRequest) Option {
	return func(a *Action) error {
		a.DeleteUserRequest = in
		return nil
	}
}

// WithUpdateUserRequest adds UpdateUserRequest to an Action struct.
func WithUpdateUserRequest(in *v1.UpdateUserRequest) Option {
	return func(a *Action) error {
		a.UpdateUserRequest = in
		return nil
	}
}

// WithDeactivateSOLRequest adds a DeactivateSOLRequest to the Action.
func WithDeactivateSOLRequest(in *v1.DeactivateSOLRequest) Option {
	return func(a *Action) error {
		a.DeactivateSOLRequest = in
		return nil
	}
}

// WithResetRequest adds ResetRequest to an Action struct.
func WithResetRequest(in *v1.ResetRequest) Option {
	return func(a *Action) error {
		a.ResetBMCRequest = in
		return nil
	}
}

// WithSkipRedfishVersions sets the Redfish versions to skip in the Action struct.
func WithSkipRedfishVersions(versions []string) Option {
	return func(a *Action) error {
		a.SkipRedfishVersions = versions
		return nil
	}
}

// NewBMC returns an oob.BMC interface.
func NewBMC(opts ...Option) (oob.BMC, error) {
	a := &Action{}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// NewBMCResetter returns an oob.BMCResetter interface.
func NewBMCResetter(opts ...Option) (*Action, error) {
	a := &Action{}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (m Action) closeConnections(ctx context.Context, connections []oob.BMC) {
	for _, conn := range connections {
		if r, ok := conn.(common.Connection); ok {
			defer r.Close(ctx) //nolint:revive // defer in a loop is OK here, as loop length is limited
		}
	}
}

// setupConnection connects to the BMC, returning BMC management methods.
func (m Action) setupConnection(ctx context.Context, user, password, host string, creds *v1.UserCreds) ([]oob.BMC, error) {
	connections := map[string]interface{}{
		"bmclibv2": &bmclibv2UserManagement{
			user:                user,
			password:            password,
			host:                host,
			log:                 m.Log,
			creds:               creds,
			skipRedfishVersions: m.SkipRedfishVersions,
		},
		"bmclibv1": &bmclibUserManagement{
			user:     user,
			password: password,
			host:     host,
			log:      m.Log,
			creds:    creds,
		},
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))

	m.SendStatusMessage("connecting to BMC")
	successfulConnections, err := common.EstablishConnections(ctx, connections)
	if err != nil {
		m.SendStatusMessage("connecting to BMC failed")
		span.SetStatus(codes.Error, "connecting to BMC failed")
		return nil, err
	}

	m.SendStatusMessage("connected to BMC")

	var actions []oob.BMC
	for _, elem := range successfulConnections {
		conn := connections[elem]
		// the connection is closed here only for bmclibv1 since bmclibv2 behaves differently
		if r, ok := conn.(common.Connection); ok && elem == "bmclibv1" {
			defer r.Close(ctx) //nolint:revive // defer in a loop is OK here, as loop length is limited
		}

		if r, ok := conn.(oob.BMC); ok {
			actions = append(actions, r)
		}
	}

	return actions, nil
}

// noteError will send a GRPC status message, log it, and update
// the status of the provided tracing span.
func (m Action) noteError(message string, span trace.Span) {
	m.SendStatusMessage(message)
	m.Log.Info(message)
	if span != nil {
		span.SetStatus(codes.Error, message)
	}
}

// CreateUser functionality for machines.
func (m Action) CreateUser(ctx context.Context) error {
	timer := prometheus.NewTimer(metrics.ActionDuration.With(prometheus.Labels{"service": "bmc", "action": "create_user"}))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "client.CreateUser")
	defer span.End()

	host, user, password, err := m.ParseAuth(m.CreateUserRequest.Authn)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))

	creds := m.CreateUserRequest.GetUserCreds()
	status := fmt.Sprintf("creating user %q", creds.GetUsername())
	m.SendStatusMessage(status)

	actions, err := m.setupConnection(ctx, user, password, host, creds)
	if err != nil {
		// setupConnection is responsible for sending a status message and updating the span
		return err
	}

	defer m.closeConnections(ctx, actions)

	if err = oob.CreateUser(ctx, actions); err != nil {
		m.noteError(fmt.Sprintf("error %s: %v", status, err), span)
		return err
	}

	m.SendStatusMessage(status + " complete")
	return nil
}

// UpdateUser functionality for machines.
func (m Action) UpdateUser(ctx context.Context) error {
	timer := prometheus.NewTimer(metrics.ActionDuration.With(prometheus.Labels{"service": "bmc", "action": "update_user"}))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "client.UpdateUser")
	defer span.End()

	host, user, password, err := m.ParseAuth(m.UpdateUserRequest.Authn)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))

	creds := m.UpdateUserRequest.GetUserCreds()
	status := fmt.Sprintf("updating user %q", creds.GetUsername())
	m.SendStatusMessage(status)

	actions, err := m.setupConnection(ctx, user, password, host, creds)
	if err != nil {
		// setupConnection is responsible for sending a status message and updating the span
		return err
	}

	defer m.closeConnections(ctx, actions)

	if err = oob.UpdateUser(ctx, actions); err != nil {
		m.noteError(fmt.Sprintf("error %s: %v", status, err), span)
		return err
	}

	m.SendStatusMessage(status + " complete")
	return nil
}

// DeleteUser functionality for machines.
func (m Action) DeleteUser(ctx context.Context) error {
	timer := prometheus.NewTimer(metrics.ActionDuration.With(prometheus.Labels{"service": "bmc", "action": "delete_user"}))
	defer timer.ObserveDuration()

	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "client.DeleteUser")
	defer span.End()

	host, user, password, err := m.ParseAuth(m.DeleteUserRequest.Authn)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))

	creds := &v1.UserCreds{Username: m.DeleteUserRequest.Username}
	status := fmt.Sprintf("deleting user %q", creds.GetUsername())
	m.SendStatusMessage(status)

	actions, err := m.setupConnection(ctx, user, password, host, creds)
	if err != nil {
		// setupConnection is responsible for sending a status message and updating the span
		return err
	}

	defer m.closeConnections(ctx, actions)

	if err = oob.DeleteUser(ctx, actions); err != nil {
		m.noteError(fmt.Sprintf("error %s: %v", status, err), span)
		return err
	}

	m.SendStatusMessage(status + " complete")
	return nil
}

// BMCReset functionality for machines.
func (m Action) BMCReset(ctx context.Context, rType string) (err error) {
	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "client.BMCReset")
	defer span.End()

	host, user, password, parseErr := m.ParseAuth(m.ResetBMCRequest.Authn)
	if parseErr != nil {
		return parseErr
	}
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))
	m.SendStatusMessage("working on bmc reset")

	opts := []bmclib.Option{
		bmclib.WithLogger(m.Log),
		bmclib.WithPerProviderTimeout(common.BMCTimeoutFromCtx(ctx)),
		bmclib.WithIpmitoolPort("623"),
	}

	client := bmclib.NewClient(host, user, password, opts...)

	lookup := map[string]string{
		v1.ResetKind_RESET_KIND_COLD.String(): "cold",
		v1.ResetKind_RESET_KIND_WARM.String(): "warm",
	}
	rLookup, ok := lookup[rType]
	if !ok {
		span.SetStatus(codes.Error, "unknown reset request")
		return &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "unknown reset request",
		}
	}
	err = client.Open(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "Permission Denied: "+err.Error())
		return &repository.Error{
			Code:    v1.Code_value["PERMISSION_DENIED"],
			Message: err.Error(),
		}
	}
	log := m.Log.WithValues("resetType", rLookup, "host", host, "user", user)
	defer func() {
		client.Close(ctx)
		log.Info("closed connections", logMetadata(client.GetMetadata())...)
	}()
	log.Info("connected to BMC", logMetadata(client.GetMetadata())...)
	m.SendStatusMessage("connected to BMC")

	ok, err = client.ResetBMC(ctx, rLookup)
	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	if err != nil {
		span.SetStatus(codes.Error, "failed to reset BMC: "+err.Error())
		log.Error(err, "failed to reset BMC")
	} else if !ok {
		err = fmt.Errorf("reset failed")
	}
	if err != nil {
		span.SetStatus(codes.Error, "failed to reset BMC: "+err.Error())
		m.SendStatusMessage(fmt.Sprintf("failed to %v reset BMC", rLookup))
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	log.Info(fmt.Sprintf("%v reset complete", rLookup))
	m.SendStatusMessage(fmt.Sprintf("%v bmc reset complete", rLookup))

	return nil
}

// DeactivateSOL deactivates a serial-over-LAN session on the device.
func (m Action) DeactivateSOL(ctx context.Context) error {
	tracer := otel.Tracer("pbnj")
	ctx, span := tracer.Start(ctx, "client.DeactivateSOL")
	defer span.End()

	host, user, password, parseErr := m.ParseAuth(m.DeactivateSOLRequest.Authn)
	if parseErr != nil {
		return parseErr
	}
	span.SetAttributes(attribute.String("bmc.host", host), attribute.String("bmc.username", user))
	m.SendStatusMessage("working on SOL session deactivation")

	opts := []bmclib.Option{
		bmclib.WithLogger(m.Log),
		bmclib.WithPerProviderTimeout(common.BMCTimeoutFromCtx(ctx)),
		bmclib.WithIpmitoolPort("623"),
	}

	client := bmclib.NewClient(host, user, password, opts...)

	if err := client.Open(ctx); err != nil {
		span.SetStatus(codes.Error, "permission denied: "+err.Error())
		return &repository.Error{
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

	err := client.DeactivateSOL(ctx)
	log = m.Log.WithValues(logMetadata(client.GetMetadata())...)
	if err != nil {
		span.SetStatus(codes.Error, "failed to deactivate SOL session: "+err.Error())
		log.Error(err, "failed to deactivate SOL session")
		m.SendStatusMessage("failed to deactivate SOL session")
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	log.Info("SOL deactivation complete")
	m.SendStatusMessage("SOL deactivation complete")

	return nil
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
