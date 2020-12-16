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
			bmcSvc := BmcService{Log: log}
			response, err := bmcSvc.NetworkSource(ctx, testCase.req)

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

func TestReset(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.ResetRequest
		message     string
		expectedErr error
	}{
		{
			name: "status good; direct auth",
			req: &v1.ResetRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "127.0.1.1",
							},
							Username: "ADMIN",
							Password: "ADMIN",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				ResetKind: v1.ResetKind_RESET_KIND_COLD,
			},
			message: "good",
		},
		{
			name:        "validation failure",
			req:         &v1.ResetRequest{Authn: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{}}}},
			message:     "",
			expectedErr: errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Username: value '' must not be an empty string"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)
			bmcSvc := BmcService{Log: log, TaskRunner: taskRunner}
			response, err := bmcSvc.Reset(ctx, testCase.req)

			t.Log("Got : ", response)
			if err != nil {
				diff := cmp.Diff(testCase.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else {
				g.Expect(response.TaskId).Should(gomega.HaveLen(20))
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.CreateUserRequest
		expectedErr error
	}{
		{
			name: "status good; direct auth",
			req: &v1.CreateUserRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "127.0.1.1",
							},
							Username: "ADMIN",
							Password: "ADMIN",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				UserCreds: &v1.UserCreds{
					Username: "admin",
					Password: "admin",
					UserRole: 0,
				},
			},
		},
		{
			name:        "validation failure",
			req:         &v1.CreateUserRequest{Authn: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{}}}},
			expectedErr: errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Username: value '' must not be an empty string"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)
			bmcSvc := BmcService{
				Log:        log,
				TaskRunner: taskRunner,
			}
			response, err := bmcSvc.CreateUser(ctx, testCase.req)

			t.Log("Got : ", response)
			if err != nil {
				diff := cmp.Diff(testCase.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}

			} else {
				g.Expect(response.TaskId).Should(gomega.HaveLen(20))
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.UpdateUserRequest
		expectedErr error
	}{
		{
			name: "status good; direct auth",
			req: &v1.UpdateUserRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "127.0.1.1",
							},
							Username: "ADMIN",
							Password: "ADMIN",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				UserCreds: &v1.UserCreds{
					Username: "admin",
					Password: "admin",
					UserRole: 0,
				},
			},
		},
		{
			name:        "validation failure",
			req:         &v1.UpdateUserRequest{Authn: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{}}}},
			expectedErr: errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Username: value '' must not be an empty string"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)
			bmcSvc := BmcService{
				Log:        log,
				TaskRunner: taskRunner,
			}
			response, err := bmcSvc.UpdateUser(ctx, testCase.req)

			t.Log("Got : ", response)
			if err != nil {
				diff := cmp.Diff(testCase.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}

			} else {
				g.Expect(response.TaskId).Should(gomega.HaveLen(20))
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.DeleteUserRequest
		expectedErr error
	}{
		{
			name: "status good; direct auth",
			req: &v1.DeleteUserRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "127.0.1.1",
							},
							Username: "ADMIN",
							Password: "ADMIN",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				Username: "me",
			},
		},
		{
			name:        "validation failure",
			req:         &v1.DeleteUserRequest{Authn: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{}}}},
			expectedErr: errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Username: value '' must not be an empty string"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)
			bmcSvc := BmcService{
				Log:        log,
				TaskRunner: taskRunner,
			}
			response, err := bmcSvc.DeleteUser(ctx, testCase.req)

			t.Log("Got : ", response)
			if err != nil {
				diff := cmp.Diff(testCase.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}

			} else {
				g.Expect(response.TaskId).Should(gomega.HaveLen(20))
			}
		})
	}
}
