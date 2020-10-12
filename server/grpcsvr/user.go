package grpcsvr

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
)

type userService struct {
	log        logging.Logger
	taskRunner task.Task
}

func (c *userService) createUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	l := c.log.GetContextLogger(ctx)
	l.V(0).Info("creating user")

	switch in.GetAuthn().Authn.(type) {
	case *v1.Authn_ExternalAuthn:
		l.V(1).Info("using external authn")
	default:
		l.V(1).Info("using direct authn")
	}

	return &v1.CreateUserResponse{
		TaskId: "user created",
	}, nil
}

func (c *userService) deleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	l := c.log.GetContextLogger(ctx)
	l.V(0).Info("deleting user")

	switch in.GetAuthn().Authn.(type) {
	case *v1.Authn_ExternalAuthn:
		l.V(1).Info("using external authn")
	default:
		l.V(1).Info("using direct authn")
	}

	return &v1.DeleteUserResponse{
		TaskId: "user deleted",
	}, nil
}

func (c *userService) updateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	l := c.log.GetContextLogger(ctx)
	l.V(0).Info("updating user")

	switch in.GetAuthn().Authn.(type) {
	case *v1.Authn_ExternalAuthn:
		l.V(1).Info("using external authn")
	default:
		l.V(1).Info("using direct authn")
	}

	return &v1.UpdateUserResponse{
		TaskId: "user updated",
	}, nil
}
