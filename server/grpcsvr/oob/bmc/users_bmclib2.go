package bmc

import (
	"context"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type bmclibClient struct {
	log      logr.Logger
	conn     *bmclib.Client
	user     string
	password string
	host     string
	creds    *v1.UserCreds
}

func (b *bmclibClient) Connect(ctx context.Context) error {
	b.conn = bmclib.NewClient(b.host, "623", b.user, b.password, bmclib.WithLogger(b.log))
	b.conn.Registry.Drivers = b.conn.Registry.Using("redfish")
	return b.conn.Open(ctx)
}

func (b *bmclibClient) Close(ctx context.Context) {
	b.conn.Close(ctx)
}

func (b *bmclibClient) CreateUser(ctx context.Context) (err error) {
	ok, err := b.conn.CreateUser(ctx, b.creds.Username, b.creds.Password, redfishRoleString(b.creds.UserRole))
	if err != nil {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	if !ok {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "error creating user",
		}
	}
	return nil
}

func (b *bmclibClient) UpdateUser(ctx context.Context) error {
	ok, err := b.conn.UpdateUser(ctx, b.creds.Username, b.creds.Password, redfishRoleString(b.creds.UserRole))
	if err != nil {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}

	if !ok {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "error updating user",
		}
	}

	return nil
}

func (b *bmclibClient) DeleteUser(ctx context.Context) error {
	return &repository.Error{
		Code:    v1.Code_value["UNKNOWN"],
		Message: "user delete not implemented",
	}
}

// redfishRoleString returns a user account role name
func redfishRoleString(role v1.UserRole) string {
	var r string
	switch role.String() {
	case "USER_ROLE_USER":
		r = "Operator"
	default:
		r = "Administrator"
	}
	return r
}
