package grpcsvr

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/onsi/gomega"
	"github.com/tinkerbell/pbnj/cmd/zaplog"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
)

func TestUserCreate(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.CreateUserRequest
		message     string
		expectedErr bool
	}{
		{
			name: "status good; direct auth",
			req: &v1.CreateUserRequest{
				Authn: &v1.Authn{
					Authn: nil,
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				UserCreds: &v1.UserCreds{
					Username: "admin",
					Password: "admin",
				},
			},
			message:     "user created",
			expectedErr: false,
		},
		{
			name: "status good; external auth",
			req: &v1.CreateUserRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_ExternalAuthn{
						ExternalAuthn: &v1.ExternalAuthn{
							Host: &v1.Host{
								Host: "10.1.1.1",
							},
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				UserCreds: &v1.UserCreds{
					Username: "admin",
					Password: "admin",
				},
			},
			message:     "user created",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)

			ctx := context.Background()

			logger, zapLogger, _ := zaplog.RegisterLogger()
			ctx = ctxzap.ToContext(ctx, zapLogger)
			userSvc := userService{
				log: logger,
			}
			response, err := userSvc.createUser(ctx, testCase.req)

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

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.DeleteUserRequest
		message     string
		expectedErr bool
	}{
		{
			name: "status good; direct auth",
			req: &v1.DeleteUserRequest{
				Authn: &v1.Authn{
					Authn: nil,
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				Username: "admin",
			},
			message:     "user deleted",
			expectedErr: false,
		},
		{
			name: "status good; external auth",
			req: &v1.DeleteUserRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_ExternalAuthn{
						ExternalAuthn: &v1.ExternalAuthn{
							Host: &v1.Host{
								Host: "10.1.1.1",
							},
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				Username: "admin",
			},
			message:     "user deleted",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)

			ctx := context.Background()

			logger, zapLogger, _ := zaplog.RegisterLogger()
			ctx = ctxzap.ToContext(ctx, zapLogger)
			userSvc := userService{
				log: logger,
			}
			response, err := userSvc.deleteUser(ctx, testCase.req)

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

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.UpdateUserRequest
		message     string
		expectedErr bool
	}{
		{
			name: "status good; direct auth",
			req: &v1.UpdateUserRequest{
				Authn: &v1.Authn{
					Authn: nil,
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				UserCreds: &v1.UserCreds{
					Username: "admin",
					Password: "admin",
				},
			},
			message:     "user updated",
			expectedErr: false,
		},
		{
			name: "status good; external auth",
			req: &v1.UpdateUserRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_ExternalAuthn{
						ExternalAuthn: &v1.ExternalAuthn{
							Host: &v1.Host{
								Host: "10.1.1.1",
							},
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				UserCreds: &v1.UserCreds{
					Username: "admin",
					Password: "admin",
				},
			},
			message:     "user updated",
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)

			ctx := context.Background()

			logger, zapLogger, _ := zaplog.RegisterLogger()
			ctx = ctxzap.ToContext(ctx, zapLogger)
			userSvc := userService{
				log: logger,
			}
			response, err := userSvc.updateUser(ctx, testCase.req)

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
