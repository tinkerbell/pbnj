package bmc

import (
	"context"
	"fmt"

	"github.com/bmc-toolbox/bmclib"
	"github.com/bmc-toolbox/bmclib/bmc"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	common "github.com/tinkerbell/pbnj/grpc/oob"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// Action for making bmc actions on BMCs, implements oob.User interface.
type Action struct {
	common.Accessory
	CreateUserRequest *v1.CreateUserRequest
	DeleteUserRequest *v1.DeleteUserRequest
	UpdateUserRequest *v1.UpdateUserRequest
	ResetBMCRequest   *v1.ResetRequest
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

// WithResetRequest adds ResetRequest to an Action struct.
func WithResetRequest(in *v1.ResetRequest) Option {
	return func(a *Action) error {
		a.ResetBMCRequest = in
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

// setupConnection connects to the BMC, returning BMC management methods.
func (m Action) setupConnection(ctx context.Context, u *bmclibUserManagement) ([]oob.BMC, error) {
	connections := map[string]interface{}{"bmclib": u}

	m.SendStatusMessage("connecting to BMC")
	successfulConnections, err := common.EstablishConnections(ctx, connections)
	if err != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return nil, err
	}

	m.SendStatusMessage("connected to BMC")

	var actions []oob.BMC
	for _, elem := range successfulConnections {
		conn := connections[elem]

		if r, ok := conn.(common.Connection); ok {
			defer r.Close(ctx) // nolint:revive // defer in a loop is OK here, as loop length is limited
		}

		if r, ok := conn.(oob.BMC); ok {
			actions = append(actions, r)
		}
	}

	return actions, nil
}

// CreateUser functionality for machines.
func (m Action) CreateUser(ctx context.Context) error {
	timer := prometheus.NewTimer(metrics.ActionDuration.With(prometheus.Labels{"service": "bmc", "action": "create_user"}))
	defer timer.ObserveDuration()

	host, user, password, err := m.ParseAuth(m.CreateUserRequest.Authn)
	if err != nil {
		return err
	}

	creds := m.CreateUserRequest.GetUserCreds()
	status := fmt.Sprintf("updating user %q", creds.GetUsername())
	m.SendStatusMessage(status)

	actions, err := m.setupConnection(ctx, &bmclibUserManagement{user: user, password: password, host: host, log: m.Log, creds: creds})
	if err != nil {
		m.SendStatusMessage("connection setup failed")
		return err
	}

	err = oob.CreateUser(ctx, actions)
	if err != nil {
		m.SendStatusMessage(fmt.Sprintf("error %s: %v", status, err))
		m.Log.Info(fmt.Sprintf("error %s: %v", status, err))
		return err
	}

	m.SendStatusMessage(status + " complete")
	return nil
}

// UpdateUser functionality for machines.
func (m Action) UpdateUser(ctx context.Context) error {
	timer := prometheus.NewTimer(metrics.ActionDuration.With(prometheus.Labels{"service": "bmc", "action": "update_user"}))
	defer timer.ObserveDuration()

	host, user, password, err := m.ParseAuth(m.UpdateUserRequest.Authn)
	if err != nil {
		return err
	}

	creds := m.UpdateUserRequest.GetUserCreds()
	status := fmt.Sprintf("updating user %q", creds.GetUsername())
	m.SendStatusMessage(status)

	actions, err := m.setupConnection(ctx, &bmclibUserManagement{user: user, password: password, host: host, log: m.Log, creds: creds})
	if err != nil {
		m.SendStatusMessage("connection setup failed")
		return err
	}

	if err = oob.UpdateUser(ctx, actions); err != nil {
		m.SendStatusMessage(fmt.Sprintf("error %s: %v", status, err))
		m.Log.Info(fmt.Sprintf("error %s: %v", status, err))
		return err
	}

	m.SendStatusMessage(status + " complete")
	return nil
}

// DeleteUser functionality for machines.
func (m Action) DeleteUser(ctx context.Context) error {
	timer := prometheus.NewTimer(metrics.ActionDuration.With(prometheus.Labels{"service": "bmc", "action": "Delete_user"}))
	defer timer.ObserveDuration()

	host, user, password, err := m.ParseAuth(m.DeleteUserRequest.Authn)
	if err != nil {
		return err
	}

	creds := &v1.UserCreds{Username: m.DeleteUserRequest.Username}
	status := fmt.Sprintf("deleting user %q", creds.GetUsername())
	m.SendStatusMessage(status)

	actions, err := m.setupConnection(ctx, &bmclibUserManagement{user: user, password: password, host: host, log: m.Log, creds: creds})
	if err != nil {
		m.SendStatusMessage("connectiion setup failed")
		return err
	}

	if err = oob.DeleteUser(ctx, actions); err != nil {
		m.SendStatusMessage(fmt.Sprintf("error %s: %v", status, err))
		m.Log.Info(fmt.Sprintf("error %s: %v", status, err))
		return err
	}

	m.SendStatusMessage(status + " complete")
	return nil
}

// BMCReset functionality for machines.
func (m Action) BMCReset(ctx context.Context, rType string) (err error) {
	host, user, password, parseErr := m.ParseAuth(m.ResetBMCRequest.Authn)
	if parseErr != nil {
		return parseErr
	}
	m.SendStatusMessage("working on bmc reset")
	client := bmclib.NewClient(host, "623", user, password, bmclib.WithLogger(m.Log))

	lookup := map[string]string{
		v1.ResetKind_RESET_KIND_COLD.String(): "cold",
		v1.ResetKind_RESET_KIND_WARM.String(): "warm",
	}
	rLookup, ok := lookup[rType]
	if !ok {
		return &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "unknown reset request",
		}
	}
	err = client.Open(ctx)
	if err != nil {
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
		log.Error(err, "failed to reset BMC")
	} else if !ok {
		err = fmt.Errorf("reset failed")
	}
	if err != nil {
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

func logMetadata(md bmc.Metadata) []interface{} {
	kvs := []interface{}{
		"ProvidersAttempted", md.ProvidersAttempted,
		"SuccessfulOpenConns", md.SuccessfulOpenConns,
		"SuccessfulCloseConns", md.SuccessfulCloseConns,
		"SuccessfulProvider", md.SuccessfulProvider,
	}

	return kvs
}
