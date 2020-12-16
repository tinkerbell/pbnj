package rpc

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/onsi/gomega"
	"github.com/packethost/pkg/log/logr"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/zaplog"
	"github.com/tinkerbell/pbnj/server/grpcsvr/persistence"
	"github.com/tinkerbell/pbnj/server/grpcsvr/taskrunner"
)

var (
	log        logging.Logger
	ctx        context.Context
	taskRunner *taskrunner.Runner
	bmcService BmcService
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	l, zapLogger, _ := logr.NewPacketLogr()
	log = zaplog.RegisterLogger(l)
	ctx = ctxzap.ToContext(ctx, zapLogger)
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	repo := &persistence.GoKV{
		Store: s,
		Ctx:   ctx,
	}

	taskRunner = &taskrunner.Runner{
		Repository: repo,
		Ctx:        ctx,
		Log:        log,
	}
	bmcService = BmcService{
		Log:                    log,
		TaskRunner:             taskRunner,
		UnimplementedBMCServer: v1.UnimplementedBMCServer{},
	}
	os.Exit(m.Run())
}

func TestConfigNetworkSource(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.NetworkSourceRequest
		message     string
		expectedErr bool
	}{
		{
			name: "status good",
			req: &v1.NetworkSourceRequest{
				Authn: &v1.Authn{
					Authn: nil,
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				NetworkSource: 0,
			},
			message:     "good",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)
			response, err := bmcService.NetworkSource(ctx, testCase.req)
			t.Log("Got : ", response)

			if testCase.expectedErr {
				g.Expect(response).ToNot(gomega.BeNil(), "Result should be nil")
				g.Expect(err).ToNot(gomega.BeNil(), "Result should be nil")
			} else {
				g.Expect(response.TaskId).To(gomega.Equal(testCase.message))
			}
		})
	}
}

func newResetRequest(authErr bool) *v1.ResetRequest {
	var auth *v1.DirectAuthn
	if authErr {
		auth = &v1.DirectAuthn{
			Host: &v1.Host{
				Host: "",
			},
			Username: "ADMIN",
			Password: "ADMIN",
		}
	} else {
		auth = &v1.DirectAuthn{
			Host: &v1.Host{
				Host: "127.0.0.1",
			},
			Username: "ADMIN",
			Password: "ADMIN",
		}
	}
	return &v1.ResetRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: auth,
			},
		},
		Vendor: &v1.Vendor{
			Name: "local",
		},
		ResetKind: 0,
	}
}
func TestReset(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		in          *v1.ResetRequest
		out         *v1.ResetResponse
	}{
		{"success", nil, newResetRequest(false), &v1.ResetResponse{TaskId: ""}},
		{"missing auth err", errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Host.Host: value '' must not be an empty string"), newResetRequest(true), nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			response, err := bmcService.Reset(ctx, tc.in)
			if err != nil {
				diff := cmp.Diff(tc.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else {
				if response.TaskId == "" {
					t.Fatal("expected taskId, got:", response.TaskId)
				}
			}
		})
	}
}

func newCreateUserRequest(authErr bool) *v1.CreateUserRequest {
	var auth *v1.DirectAuthn
	if authErr {
		auth = &v1.DirectAuthn{
			Host: &v1.Host{
				Host: "",
			},
			Username: "ADMIN",
			Password: "ADMIN",
		}
	} else {
		auth = &v1.DirectAuthn{
			Host: &v1.Host{
				Host: "127.0.0.1",
			},
			Username: "ADMIN",
			Password: "ADMIN",
		}
	}
	return &v1.CreateUserRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: auth,
			},
		},
		Vendor: &v1.Vendor{
			Name: "local",
		},
		UserCreds: &v1.UserCreds{
			Username: "",
			Password: "",
			UserRole: 0,
		},
	}
}
func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		in          *v1.CreateUserRequest
		out         *v1.CreateUserResponse
	}{
		{"success", nil, newCreateUserRequest(false), &v1.CreateUserResponse{TaskId: ""}},
		{"missing auth err", errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Host.Host: value '' must not be an empty string"), newCreateUserRequest(true), nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			response, err := bmcService.CreateUser(ctx, tc.in)
			if err != nil {
				diff := cmp.Diff(tc.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else {
				if response.TaskId == "" {
					t.Fatal("expected taskId, got:", response.TaskId)
				}
			}
		})
	}
}

func newUpdateUserRequest(authErr bool) *v1.UpdateUserRequest {
	var auth *v1.DirectAuthn
	if authErr {
		auth = &v1.DirectAuthn{
			Host: &v1.Host{
				Host: "",
			},
			Username: "ADMIN",
			Password: "ADMIN",
		}
	} else {
		auth = &v1.DirectAuthn{
			Host: &v1.Host{
				Host: "127.0.0.1",
			},
			Username: "ADMIN",
			Password: "ADMIN",
		}
	}
	return &v1.UpdateUserRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: auth,
			},
		},
		Vendor: &v1.Vendor{
			Name: "local",
		},
		UserCreds: &v1.UserCreds{
			Username: "",
			Password: "",
			UserRole: 0,
		},
	}
}
func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		in          *v1.UpdateUserRequest
		out         *v1.UpdateUserResponse
	}{
		{"success", nil, newUpdateUserRequest(false), &v1.UpdateUserResponse{TaskId: ""}},
		{"missing auth err", errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Host.Host: value '' must not be an empty string"), newUpdateUserRequest(true), nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			response, err := bmcService.UpdateUser(ctx, tc.in)
			if err != nil {
				diff := cmp.Diff(tc.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else {
				if response.TaskId == "" {
					t.Fatal("expected taskId, got:", response.TaskId)
				}
			}
		})
	}
}

func newDeleteUserRequest(authErr bool) *v1.DeleteUserRequest {
	var auth *v1.DirectAuthn
	if authErr {
		auth = &v1.DirectAuthn{
			Host: &v1.Host{
				Host: "",
			},
			Username: "ADMIN",
			Password: "ADMIN",
		}
	} else {
		auth = &v1.DirectAuthn{
			Host: &v1.Host{
				Host: "127.0.0.1",
			},
			Username: "ADMIN",
			Password: "ADMIN",
		}
	}
	return &v1.DeleteUserRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: auth,
			},
		},
		Vendor: &v1.Vendor{
			Name: "local",
		},
		Username: "blah",
	}
}
func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		in          *v1.DeleteUserRequest
		out         *v1.UpdateUserResponse
	}{
		{"success", nil, newDeleteUserRequest(false), &v1.UpdateUserResponse{TaskId: ""}},
		{"missing auth err", errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Host.Host: value '' must not be an empty string"), newDeleteUserRequest(true), nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			response, err := bmcService.DeleteUser(ctx, tc.in)
			if err != nil {
				diff := cmp.Diff(tc.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else {
				if response.TaskId == "" {
					t.Fatal("expected taskId, got:", response.TaskId)
				}
			}
		})
	}
}
