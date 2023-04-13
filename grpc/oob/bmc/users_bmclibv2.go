package bmc

import (
	"context"
	"fmt"

	"github.com/bmc-toolbox/bmclib/v2"
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	common "github.com/tinkerbell/pbnj/grpc/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// bmclibv2UserManagement wraps attributes to manage user accounts with bmclib v2.
type bmclibv2UserManagement struct {
	conn     *bmclib.Client
	host     string
	log      logr.Logger
	password string
	user     string
	creds    *v1.UserCreds
	// skipRedfishVersions is a list of Redfish versions to be ignored,
	//
	// When running an action on a BMC, PBnJ will pass the value of the skipRedfishVersions to bmclib
	// which will then ignore the Redfish endpoint completely on BMCs running the given Redfish versions,
	// and will proceed to attempt other drivers like - IPMI/SSH/Vendor API instead.
	//
	// for more information see https://github.com/bmc-toolbox/bmclib#bmc-connections
	skipRedfishVersions []string
}

// Connect sets up the BMC client connection.
func (b *bmclibv2UserManagement) Connect(ctx context.Context) error {
	var errMsg repository.Error

	opts := []bmclib.Option{
		bmclib.WithLogger(b.log),
		bmclib.WithPerProviderTimeout(common.BMCTimeoutFromCtx(ctx)),
	}

	if len(b.skipRedfishVersions) > 0 {
		opts = append(opts, bmclib.WithRedfishVersionsNotCompatible(b.skipRedfishVersions))
	}

	client := bmclib.NewClient(b.host, "", b.user, b.password, opts...)
	client.Registry.Drivers = client.Registry.FilterForCompatible(ctx)

	if err := client.Open(ctx); err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return &errMsg
	}

	b.conn = client

	return nil
}

// Close closes the BMC client connection.
func (b *bmclibv2UserManagement) Close(ctx context.Context) {
	b.conn.Close(ctx)
}

// CreateUser creates a user account.
func (b *bmclibv2UserManagement) CreateUser(ctx context.Context) error {
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
func (b *bmclibv2UserManagement) UpdateUser(ctx context.Context) error {
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
func (b *bmclibv2UserManagement) DeleteUser(ctx context.Context) error {
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
