package bmc

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
	bmc "github.com/tinkerbell/pbnj/server/grpcsvr/oob"
)

// Action for making bmc actions on BMCs, implements oob.User interface
type Action struct {
	bmc.Accessory
	CreateUserRequest *v1.CreateUserRequest
	DeleteUserRequest *v1.DeleteUserRequest
	UpdateUserRequest *v1.UpdateUserRequest
}

type userConnection struct {
	bmc.ConnectionDetails
	user user
}

// the power interface allows us to abstract these functions
// between different libraries and BMC connections
// like ipmi, racadm, redfish, etc
type user interface {
	bmc.Connection
	create(context.Context) repository.Error
	update(context.Context) repository.Error
	delete(context.Context) repository.Error
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

// CreateUser functionality for machines
func (m Action) CreateUser(ctx context.Context) (errMsg repository.Error) {
	host, user, password, errMsg := m.ParseAuth(m.CreateUserRequest.Authn)
	if errMsg.Message != "" {
		return errMsg
	}
	creds := m.CreateUserRequest.GetUserCreds()
	base := "creating user: " + creds.GetUsername()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := []userConnection{
		{
			ConnectionDetails: bmc.ConnectionDetails{Name: "bmclib"},
			user: &bmclibUserManagement{
				user:     user,
				password: password,
				host:     host,
				creds:    creds,
			},
		},
	}

	var connected bool
	m.SendStatusMessage("connecting to BMC")
	for index := range connections {
		connections[index].Err = connections[index].user.Connect(ctx)
		if connections[index].Err.Message == "" {
			connections[index].Connected = true
			defer connections[index].user.Close(ctx)
			connected = true
		}
	}
	m.Log.V(1).Info("connections", "connections", fmt.Sprintf("%+v", connections))
	if !connected {
		m.SendStatusMessage("connecting to BMC failed")
		var combinedErrs []string
		for _, connection := range connections {
			combinedErrs = append(combinedErrs, connection.Err.Message)
		}
		msg := "could not connect"
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = msg
		errMsg.Details = append(errMsg.Details, combinedErrs...)
		m.Log.V(0).Info(msg, "error", combinedErrs)
		return errMsg
	}
	m.SendStatusMessage("connected to BMC")

	for index := range connections {
		if connections[index].Connected {
			m.Log.V(0).Info("trying", "name", connections[index].Name)
			errMsg = connections[index].user.create(ctx)
			if errMsg.Message == "" {
				m.Log.V(0).Info("action implemented by", "implementer", connections[index].Name)
				break
			}
		}
	}

	if errMsg.Message != "" {
		m.SendStatusMessage("error with " + base + ": " + errMsg.Message)
		m.Log.V(0).Info("error with "+base, "error", errMsg.Message)
	}
	m.SendStatusMessage(base + " complete")
	return errMsg //nolint
}

// UpdateUser functionality for machines
func (m Action) UpdateUser(ctx context.Context) (errMsg repository.Error) {
	host, user, password, errMsg := m.ParseAuth(m.UpdateUserRequest.Authn)
	if errMsg.Message != "" {
		return errMsg
	}
	creds := m.UpdateUserRequest.GetUserCreds()
	base := "updating user: " + creds.GetUsername()
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := []userConnection{
		{
			ConnectionDetails: bmc.ConnectionDetails{Name: "bmclib"},
			user: &bmclibUserManagement{
				user:     user,
				password: password,
				host:     host,
				creds:    creds,
			},
		},
	}

	var connected bool
	m.SendStatusMessage("connecting to BMC")
	for index := range connections {
		connections[index].Err = connections[index].user.Connect(ctx)
		if connections[index].Err.Message == "" {
			connections[index].Connected = true
			defer connections[index].user.Close(ctx)
			connected = true
		}
	}
	m.Log.V(1).Info("connections", "connections", fmt.Sprintf("%+v", connections))
	if !connected {
		m.SendStatusMessage("connecting to BMC failed")
		var combinedErrs []string
		for _, connection := range connections {
			combinedErrs = append(combinedErrs, connection.Err.Message)
		}
		msg := "could not connect"
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = msg
		errMsg.Details = append(errMsg.Details, combinedErrs...)
		m.Log.V(0).Info(msg, "error", combinedErrs)
		return errMsg
	}
	m.SendStatusMessage("connected to BMC")

	for index := range connections {
		if connections[index].Connected {
			m.Log.V(0).Info("trying", "name", connections[index].Name)
			errMsg = connections[index].user.update(ctx)
			if errMsg.Message == "" {
				m.Log.V(0).Info("action implemented by", "implementer", connections[index].Name)
				break
			}
		}
	}

	if errMsg.Message != "" {
		m.SendStatusMessage("error with " + base + ": " + errMsg.Message)
		m.Log.V(0).Info("error with "+base, "error", errMsg.Message)
	}
	m.SendStatusMessage(base + " complete")
	return errMsg //nolint
}

// DeleteUser functionality for machines
func (m Action) DeleteUser(ctx context.Context) (errMsg repository.Error) {
	host, user, password, errMsg := m.ParseAuth(m.DeleteUserRequest.Authn)
	if errMsg.Message != "" {
		return errMsg
	}
	base := "deleting user: " + m.DeleteUserRequest.Username
	msg := "working on " + base
	m.SendStatusMessage(msg)

	connections := []userConnection{
		{
			ConnectionDetails: bmc.ConnectionDetails{Name: "bmclib"},
			user: &bmclibUserManagement{
				user:     user,
				password: password,
				host:     host,
				creds: &v1.UserCreds{
					Username: m.DeleteUserRequest.Username,
				},
			},
		},
	}

	var connected bool
	m.SendStatusMessage("connecting to BMC")
	for index := range connections {
		connections[index].Err = connections[index].user.Connect(ctx)
		if connections[index].Err.Message == "" {
			connections[index].Connected = true
			defer connections[index].user.Close(ctx)
			connected = true
		}
	}
	m.Log.V(1).Info("connections", "connections", fmt.Sprintf("%+v", connections))
	if !connected {
		m.SendStatusMessage("connecting to BMC failed")
		var combinedErrs []string
		for _, connection := range connections {
			combinedErrs = append(combinedErrs, connection.Err.Message)
		}
		msg := "could not connect"
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = msg
		errMsg.Details = append(errMsg.Details, combinedErrs...)
		m.Log.V(0).Info(msg, "error", combinedErrs)
		return errMsg
	}
	m.SendStatusMessage("connected to BMC")

	for index := range connections {
		if connections[index].Connected {
			m.Log.V(0).Info("trying", "name", connections[index].Name)
			errMsg = connections[index].user.delete(ctx)
			if errMsg.Message == "" {
				m.Log.V(0).Info("action implemented by", "implementer", connections[index].Name)
				break
			}
		}
	}

	if errMsg.Message != "" {
		m.SendStatusMessage("error with " + base + ": " + errMsg.Message)
		m.Log.V(0).Info("error with "+base, "error", errMsg.Message)
	}
	m.SendStatusMessage(base + " complete")
	return errMsg //nolint
}
