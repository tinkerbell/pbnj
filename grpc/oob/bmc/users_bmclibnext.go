package bmc

import (
	"context"
	"fmt"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// bmclibNextUserManagement wraps attributes to manage user accounts.
type bmclibNextUserManagement struct {
	conn     *bmclib.Client
	host     string
	log      logr.Logger
	password string
	user     string
	creds    *v1.UserCreds
}

// Connect sets up the BMC client connection.
func (b *bmclibNextUserManagement) Connect(ctx context.Context) error {
	var errMsg repository.Error

	client := bmclib.NewClient(b.host, "", b.user, b.password, bmclib.WithLogger(b.log))
	client.Registry.Drivers = client.Registry.Using("redfish")

	if err := client.Open(ctx); err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return &errMsg
	}

	b.conn = client

	return nil
}

// Close closes the BMC client connection.
func (b *bmclibNextUserManagement) Close(ctx context.Context) {
	b.conn.Close(ctx)
}

// CreateUser creates a user account.
func (b *bmclibNextUserManagement) CreateUser(ctx context.Context) error {
	ok, err := b.conn.CreateUser(ctx, b.creds.Username, b.creds.Password, userRole(b.creds.UserRole))
	if !ok && err == nil {
		err = fmt.Errorf("error creating user")
	}

	if err != nil {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	return nil
}

// UpdateUser updates a user account.
func (b *bmclibNextUserManagement) UpdateUser(ctx context.Context) error {
	ok, err := b.conn.UpdateUser(ctx, b.creds.Username, b.creds.Password, userRole(b.creds.UserRole))
	if !ok && err == nil {
		err = fmt.Errorf("error updating user")
	}

	if err != nil {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return nil
}

// DeleteUser deletes a user account.
func (b *bmclibNextUserManagement) DeleteUser(ctx context.Context) error {
	ok, err := b.conn.DeleteUser(ctx, b.creds.Username)
	if !ok && err == nil {
		err = fmt.Errorf("error deleting user")
	}

	if err != nil {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	return nil
}

// userRole returns the Redfish equivalent role for a v1.UserRole string.
func userRole(role v1.UserRole) string {
	var r string
	switch role.String() {
	case "USER_ROLE_USER":
		r = "Operator"
	default:
		r = "Administrator"
	}
	return r
}
