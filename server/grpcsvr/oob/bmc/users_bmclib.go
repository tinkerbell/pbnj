package bmc

import (
	"context"

	"github.com/bmc-toolbox/bmclib/cfgresources"
	"github.com/bmc-toolbox/bmclib/devices"
	"github.com/bmc-toolbox/bmclib/discover"
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type bmclibUserManagement struct {
	conn     devices.Bmc
	host     string
	log      logr.Logger
	password string
	user     string
	creds    *v1.UserCreds
}

func (b *bmclibUserManagement) Connect(ctx context.Context) repository.Error {
	var errMsg repository.Error
	connection, err := discover.ScanAndConnect(b.host, b.user, b.password, discover.WithLogger(b.log))
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return errMsg //nolint
	}
	switch conn := connection.(type) {
	case devices.Bmc:
		b.conn = conn
	default:
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "Unknown device"
		return errMsg //nolint
	}
	return errMsg //nolint
}

func (b *bmclibUserManagement) Close(ctx context.Context) {
	b.conn.Close()
}

func (b *bmclibUserManagement) create(ctx context.Context) repository.Error {
	var errMsg repository.Error
	users := []*cfgresources.User{
		{
			Name:     b.creds.Username,
			Password: b.creds.Password,
			Role:     userRoleToString(b.creds.UserRole),
			Enable:   true,
		},
	}
	err := b.conn.User(users)
	if err != nil {
		return repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return errMsg
}

func (b *bmclibUserManagement) update(ctx context.Context) repository.Error {
	var errMsg repository.Error
	users := []*cfgresources.User{
		{
			Name:     b.creds.Username,
			Password: b.creds.Password,
			Role:     userRoleToString(b.creds.UserRole),
			Enable:   true,
		},
	}
	err := b.conn.User(users)
	if err != nil {
		return repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return errMsg
}

func (b *bmclibUserManagement) delete(ctx context.Context) repository.Error {
	var errMsg repository.Error
	users := []*cfgresources.User{
		{
			Name:     b.creds.Username,
			Password: "DELETE",
			Role:     "user",
			Enable:   false,
		},
	}
	err := b.conn.User(users)
	if err != nil {
		return repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return errMsg
}

func userRoleToString(role v1.UserRole) string {
	var r string
	switch role.String() {
	case "USER_ROLE_USER":
		r = "user"
	default:
		r = "admin"
	}
	return r
}
