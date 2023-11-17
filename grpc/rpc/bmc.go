package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/rs/xid"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/oob/bmc"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
	"go.opentelemetry.io/otel/trace"
)

// BmcService for doing BMC actions.
type BmcService struct {
	// Timeout is how long a task should be run
	// before it is cancelled. This is for use in a
	// TaskRunner.Execute function that runs all BMC
	// interactions in the background.
	Timeout time.Duration
	// SkipRedfishVersions is a list of Redfish versions to be ignored,
	//
	// When running an action on a BMC, PBnJ will pass the value of the skipRedfishVersions to bmclib
	// which will then ignore the Redfish endpoint completely on BMCs running the given Redfish versions,
	// and will proceed to attempt other drivers like - IPMI/SSH/Vendor API instead.
	//
	// for more information see https://github.com/bmc-toolbox/bmclib#bmc-connections
	SkipRedfishVersions []string
	TaskRunner          task.Task
	v1.UnimplementedBMCServer
}

// NetworkSource sets the BMC network source.
func (b *BmcService) NetworkSource(_ context.Context, _ *v1.NetworkSourceRequest) (*v1.NetworkSourceResponse, error) {
	return nil, errors.New("not implemented")
}

// Reset calls a reset on a BMC.
func (b *BmcService) Reset(ctx context.Context, in *v1.ResetRequest) (*v1.ResetResponse, error) {
	l := logging.ExtractLogr(ctx)
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
		// Because this is a background task, we want to pass through the span context, but not be
		// a child context. This allows us to correctly plumb otel into the background task.
		c := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
		taskCtx, cancel := context.WithTimeout(c, b.Timeout)
		defer cancel()
		return "", t.BMCReset(taskCtx, in.ResetKind.String())
	}
	b.TaskRunner.Execute(ctx, l, "bmc reset", taskID, execFunc)

	return &v1.ResetResponse{TaskId: taskID}, nil
}

// DeactivateSOL deactivates any active SOL session on the BMC.
func (b *BmcService) DeactivateSOL(ctx context.Context, in *v1.DeactivateSOLRequest) (*v1.DeactivateSOLResponse, error) {
	l := logging.ExtractLogr(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)
	l.Info(
		"start DeactivateSOL request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
	)

	execFunc := func(s chan string) (string, error) {
		t, err := bmc.NewBMCResetter(
			bmc.WithDeactivateSOLRequest(in),
			bmc.WithLogger(l),
			bmc.WithStatusMessage(s),
		)
		if err != nil {
			return "", err
		}
		// Because this is a background task, we want to pass through the span context, but not be
		// a child context. This allows us to correctly plumb otel into the background task.
		c := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
		taskCtx, cancel := context.WithTimeout(c, b.Timeout)
		defer cancel()
		return "", t.DeactivateSOL(taskCtx)
	}
	b.TaskRunner.Execute(ctx, l, "deactivating SOL session", taskID, execFunc)

	return &v1.DeactivateSOLResponse{TaskId: taskID}, nil
}

// CreateUser sets the next boot device of a machine.
func (b *BmcService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	l := logging.ExtractLogr(ctx)
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
		// Because this is a background task, we want to pass through the span context, but not be
		// a child context. This allows us to correctly plumb otel into the background task.
		c := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
		taskCtx, cancel := context.WithTimeout(c, b.Timeout)
		defer cancel()
		return "", t.CreateUser(taskCtx)
	}
	b.TaskRunner.Execute(ctx, l, "creating user", taskID, execFunc)

	return &v1.CreateUserResponse{TaskId: taskID}, nil
}

// UpdateUser updates a users credentials on a BMC.
func (b *BmcService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	l := logging.ExtractLogr(ctx)
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
		// Because this is a background task, we want to pass through the span context, but not be
		// a child context. This allows us to correctly plumb otel into the background task.
		c := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
		taskCtx, cancel := context.WithTimeout(c, b.Timeout)
		defer cancel()
		return "", t.UpdateUser(taskCtx)
	}
	b.TaskRunner.Execute(ctx, l, "updating user", taskID, execFunc)

	return &v1.UpdateUserResponse{TaskId: taskID}, nil
}

// DeleteUser deletes a user on a BMC.
func (b *BmcService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	l := logging.ExtractLogr(ctx)
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
		// Because this is a background task, we want to pass through the span context, but not be
		// a child context. This allows us to correctly plumb otel into the background task.
		c := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(ctx))
		taskCtx, cancel := context.WithTimeout(c, b.Timeout)
		defer cancel()
		return "", t.DeleteUser(taskCtx)
	}
	b.TaskRunner.Execute(ctx, l, "deleting user", taskID, execFunc)

	return &v1.DeleteUserResponse{TaskId: taskID}, nil
}
