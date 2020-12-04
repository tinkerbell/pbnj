package bmc

import (
	"context"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
	common "github.com/tinkerbell/pbnj/server/grpcsvr/oob"
)

// Action for making bmc actions on BMCs, implements oob.User interface
type Action struct {
	common.Accessory
	CreateUserRequest *v1.CreateUserRequest
	DeleteUserRequest *v1.DeleteUserRequest
	UpdateUserRequest *v1.UpdateUserRequest
	ResetBMCRequest   *v1.ResetRequest
}

// Option to add to an Actions
type Option func(a *Action) error

// WithLogger adds a logr to an Action struct
func WithLogger(l logr.Logger) Option {
	return func(a *Action) error {
		a.Log = l
		return nil
	}
}

// WithStatusMessage adds a status message chan to an Action struct
func WithStatusMessage(s chan string) Option {
	return func(a *Action) error {
		a.StatusMessages = s
		return nil
	}
}

// WithCreateUserRequest adds CreateUserRequest to an Action struct
func WithCreateUserRequest(in *v1.CreateUserRequest) Option {
	return func(a *Action) error {
		a.CreateUserRequest = in
		return nil
	}
}

// WithDeleteUserRequest adds DeleteUserRequest to an Action struct
func WithDeleteUserRequest(in *v1.DeleteUserRequest) Option {
	return func(a *Action) error {
		a.DeleteUserRequest = in
		return nil
	}
}

// WithUpdateUserRequest adds UpdateUserRequest to an Action struct
func WithUpdateUserRequest(in *v1.UpdateUserRequest) Option {
	return func(a *Action) error {
		a.UpdateUserRequest = in
		return nil
	}
}

// WithResetRequest adds ResetRequest to an Action struct
func WithResetRequest(in *v1.ResetRequest) Option {
	return func(a *Action) error {
		a.ResetBMCRequest = in
		return nil
	}
}

// NewBMC returns an oob.BMC interface
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

// NewBMCResetter returns an oob.BMCResetter interface
func NewBMCResetter(opts ...Option) (oob.BMCResetter, error) {
	a := &Action{}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// CreateUser functionality for machines
func (m Action) CreateUser(ctx context.Context) error {
	labels := prometheus.Labels{
		"service": "bmc",
		"action":  "create_user",
	}
	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	var err error
	host, user, password, parseErr := m.ParseAuth(m.CreateUserRequest.Authn)
	if parseErr != nil {
		return parseErr
	}
	creds := m.UpdateUserRequest.GetUserCreds()
	base := "creating user: " + creds.GetUsername()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := map[string]interface{}{
		"bmclib": &bmclibUserManagement{user: user, password: password, host: host, creds: creds},
	}

	m.SendStatusMessage("connecting to BMC")
	successfulConnections, ecErr := common.EstablishConnections(ctx, connections)
	if ecErr != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return ecErr
	}
	m.SendStatusMessage("connected to BMC")

	var userAction []oob.BMC
	for _, elem := range successfulConnections {
		conn := connections[elem]
		switch r := conn.(type) {
		case oob.BMC:
			userAction = append(userAction, r)
		}
	}
	err = oob.CreateUser(ctx, userAction)
	if err != nil {
		m.SendStatusMessage("error with " + base + ": " + err.Error())
		m.Log.V(0).Info("error with "+base, "error", err.Error())
		return err
	}
	m.SendStatusMessage(base + " complete")
	return nil
}

// UpdateUser functionality for machines
func (m Action) UpdateUser(ctx context.Context) error {
	labels := prometheus.Labels{
		"service": "bmc",
		"action":  "update_user",
	}
	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	var err error
	host, user, password, parseErr := m.ParseAuth(m.UpdateUserRequest.Authn)
	if parseErr != nil {
		return parseErr
	}
	creds := m.UpdateUserRequest.GetUserCreds()
	base := "updating user: " + creds.GetUsername()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := map[string]interface{}{
		"bmclib": &bmclibUserManagement{
			user:     user,
			password: password,
			host:     host,
			creds:    creds,
		},
	}

	m.SendStatusMessage("connecting to BMC")
	successfulConnections, ecErr := common.EstablishConnections(ctx, connections)
	if ecErr != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return ecErr
	}
	m.SendStatusMessage("connected to BMC")

	var userAction []oob.BMC
	for _, elem := range successfulConnections {
		switch r := connections[elem].(type) {
		case oob.BMC:
			userAction = append(userAction, r)
		}
	}
	err = oob.UpdateUser(ctx, userAction)
	if err != nil {
		m.SendStatusMessage("error with " + base + ": " + err.Error())
		m.Log.V(0).Info("error with "+base, "error", err.Error())
		return err
	}
	m.SendStatusMessage(base + " complete")
	return nil
}

// DeleteUser functionality for machines
func (m Action) DeleteUser(ctx context.Context) error {
	labels := prometheus.Labels{
		"service": "bmc",
		"action":  "delete_user",
	}
	timer := prometheus.NewTimer(metrics.ActionDuration.With(labels))
	defer timer.ObserveDuration()

	var deleteErr error
	host, user, password, parseErr := m.ParseAuth(m.DeleteUserRequest.Authn)
	if parseErr != nil {
		return parseErr
	}
	base := "deleting user: " + m.DeleteUserRequest.Username
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := map[string]interface{}{
		"bmclib": &bmclibUserManagement{
			user:     user,
			password: password,
			host:     host,
			creds: &v1.UserCreds{
				Username: m.DeleteUserRequest.Username,
			},
		},
	}

	m.SendStatusMessage("connecting to BMC")
	successfulConnections, ecErr := common.EstablishConnections(ctx, connections)
	if ecErr != nil {
		m.SendStatusMessage("connecting to BMC failed")
		return ecErr
	}
	m.SendStatusMessage("connected to BMC")

	var deleteUsers []oob.BMC
	for _, elem := range successfulConnections {
		switch r := connections[elem].(type) {
		case oob.BMC:
			deleteUsers = append(deleteUsers, r)
		}
	}
	deleteErr = oob.DeleteUser(ctx, deleteUsers)
	if deleteErr != nil {
		m.SendStatusMessage("error with " + base + ": " + deleteErr.Error())
		m.Log.V(0).Info("error with "+base, "error", deleteErr.Error())
		return deleteErr
	}
	m.SendStatusMessage(base + " complete")
	return nil
}

// BMCReset functionality for machines
func (m Action) BMCReset(ctx context.Context, rType string) (err error) {
	host, user, password, parseErr := m.ParseAuth(m.ResetBMCRequest.Authn)
	if parseErr != nil {
		return parseErr
	}
	m.SendStatusMessage("working on bmc reset")
	client := bmclib.NewClient(host, "623", user, password, bmclib.WithLogger(m.Log))
	err = client.DiscoverProviders(ctx)
	if err != nil {
		m.Log.V(1).Info("error with provider discovery", "err", err.Error())
	}
	var errMsg string
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
	ok, err = client.ResetBMC(ctx, rLookup)
	if err != nil {
		errMsg = err.Error()
	} else if !ok {
		errMsg = "reset failed"
	}
	if errMsg != "" {
		err = &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: errMsg,
		}
	}
	m.SendStatusMessage("bmc reset complete")
	return err
}
