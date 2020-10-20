package rpc

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
)

// UserService for crud operations on BMC users
type UserService struct {
	Log        logging.Logger
	TaskRunner task.Task
}

// CreateUser creates a user on a BMC
func (u *UserService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	l := u.Log.GetContextLogger(ctx)
	l.V(0).Info("creating user")

	return &v1.CreateUserResponse{
		TaskId: "user created",
	}, nil
}

// DeleteUser deletes a user on a BMC
func (u *UserService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	l := u.Log.GetContextLogger(ctx)
	l.V(0).Info("deleting user")

	return &v1.DeleteUserResponse{
		TaskId: "user deleted",
	}, nil
}

// UpdateUser updates a users credentials on a BMC
func (u *UserService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	l := u.Log.GetContextLogger(ctx)
	l.V(0).Info("updating user")

	return &v1.UpdateUserResponse{
		TaskId: "user updated",
	}, nil
}
