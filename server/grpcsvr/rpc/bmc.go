package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/rs/xid"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
	"github.com/tinkerbell/pbnj/server/grpcsvr/oob/bmc"
)

// BmcService for doing BMC actions.
type BmcService struct {
	Log logging.Logger
	// Timeout is how long a task should be run
	// before it is cancelled. This is for use in a
	// TaskRunner.Execute function that runs all BMC
	// interactions in the background.
	Timeout    time.Duration
	TaskRunner task.Task
	v1.UnimplementedBMCServer
}

// NetworkSource sets the BMC network source.
func (b *BmcService) NetworkSource(_ context.Context, _ *v1.NetworkSourceRequest) (*v1.NetworkSourceResponse, error) {
	return nil, errors.New("not implemented")
}

// Reset calls a reset on a BMC.
func (b *BmcService) Reset(ctx context.Context, in *v1.ResetRequest) (*v1.ResetResponse, error) {
	l := b.Log.GetContextLogger(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)

	l.Info(
		"start Reset request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
		"resetKind", in.GetResetKind().String(),
	)

	execFunc := func(s chan string) (string, error) {
		t, err := bmc.NewBMCResetter(
			bmc.WithLogger(l),
			bmc.WithStatusMessage(s),
			bmc.WithResetRequest(in),
		)
		if err != nil {
			return "", err
		}
		taskCtx, cancel := context.WithTimeout(ctx, b.Timeout)
		// cant defer this cancel because it cancels the context before the func is run
		// cant have cancel be _ because go vet complains.
		// TODO(jacobweinstock): maybe move this context withTimeout into the TaskRunner.Execute function
		_ = cancel
		return "", t.BMCReset(taskCtx, in.ResetKind.String())
	}
	b.TaskRunner.Execute(ctx, "bmc reset", taskID, execFunc)

	return &v1.ResetResponse{TaskId: taskID}, nil
}

// CreateUser sets the next boot device of a machine.
func (b *BmcService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	// TODO figure out how not to have to do this, but still keep the logging abstraction clean?
	l := b.Log.GetContextLogger(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)

	l.Info(
		"start CreateUser request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
		"userCreds.Username", in.UserCreds.Username,
		"userCreds.UserRole", in.UserCreds.UserRole,
	)

	execFunc := func(s chan string) (string, error) {
		t, err := bmc.NewBMC(
			bmc.WithCreateUserRequest(in),
			bmc.WithLogger(l),
			bmc.WithStatusMessage(s),
		)
		if err != nil {
			return "", err
		}
		taskCtx, cancel := context.WithTimeout(context.Background(), b.Timeout)
		_ = cancel
		return "", t.CreateUser(taskCtx)
	}
	b.TaskRunner.Execute(ctx, "creating user", taskID, execFunc)

	return &v1.CreateUserResponse{TaskId: taskID}, nil
}

// UpdateUser updates a users credentials on a BMC.
func (b *BmcService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	// TODO figure out how not to have to do this, but still keep the logging abstraction clean?
	l := b.Log.GetContextLogger(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)

	l.Info(
		"start UpdateUser request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
		"userCreds.Username", in.UserCreds.Username,
		"userCreds.UserRole", in.UserCreds.UserRole,
	)

	execFunc := func(s chan string) (string, error) {
		t, err := bmc.NewBMC(
			bmc.WithUpdateUserRequest(in),
			bmc.WithLogger(l),
			bmc.WithStatusMessage(s),
		)
		if err != nil {
			return "", err
		}
		taskCtx, cancel := context.WithTimeout(context.Background(), b.Timeout)
		_ = cancel
		return "", t.UpdateUser(taskCtx)
	}
	b.TaskRunner.Execute(ctx, "updating user", taskID, execFunc)

	return &v1.UpdateUserResponse{TaskId: taskID}, nil
}

// DeleteUser deletes a user on a BMC.
func (b *BmcService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	// TODO figure out how not to have to do this, but still keep the logging abstraction clean?
	l := b.Log.GetContextLogger(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)
	l.Info(
		"start DeleteUser request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
		"userCreds.Username", in.Username,
	)

	execFunc := func(s chan string) (string, error) {
		t, err := bmc.NewBMC(
			bmc.WithDeleteUserRequest(in),
			bmc.WithLogger(l),
			bmc.WithStatusMessage(s),
		)
		if err != nil {
			return "", err
		}
		taskCtx, cancel := context.WithTimeout(context.Background(), b.Timeout)
		_ = cancel
		return "", t.DeleteUser(taskCtx)
	}
	b.TaskRunner.Execute(ctx, "deleting user", taskID, execFunc)

	return &v1.DeleteUserResponse{TaskId: taskID}, nil
}
