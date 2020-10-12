package grpcsvr

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/onsi/gomega"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	"github.com/tinkerbell/pbnj/cmd/zaplog"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
	"github.com/tinkerbell/pbnj/pkg/task"
	"github.com/tinkerbell/pbnj/server/grpcsvr/persistence"
	"github.com/tinkerbell/pbnj/server/grpcsvr/taskrunner"
)

func TestDevice(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.DeviceRequest
		message     string
		expectedErr bool
	}{
		{
			name: "status good; direct auth",
			req: &v1.DeviceRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "",
							},
							Username: "",
							Password: "",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				Persistent: false,
				EfiBoot:    false,
			},
			message:     "good",
			expectedErr: false,
		},
		{
			name: "status good; external auth",
			req: &v1.DeviceRequest{
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
				Persistent: false,
				EfiBoot:    false,
			},
			message:     "good",
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
			f := freecache.NewStore(freecache.DefaultOptions)
			s := gokv.Store(f)
			var repo repository.Actions
			repo = &persistence.GoKV{
				Store: s,
				Ctx:   ctx,
			}

			var taskRunner task.Task
			taskRunner = &taskrunner.Runner{
				Repository: repo,
				Ctx:        ctx,
				Log:        logger,
			}
			machineSvc := machineService{
				log:        logger,
				taskRunner: taskRunner,
			}
			response, err := machineSvc.device(ctx, testCase.req)

			t.Log("Got : ", response)

			if testCase.expectedErr {
				g.Expect(response).ToNot(gomega.BeNil(), "Result should be nil")
				g.Expect(err).ToNot(gomega.BeNil(), "Result should be nil")
			} else {
				g.Expect(response.TaskId).Should(gomega.HaveLen(20))
			}
		})
	}
}

func TestPower(t *testing.T) {
	// TODO make sure external auth doesnt break stuff
	testCases := []struct {
		name        string
		req         *v1.PowerRequest
		message     string
		expectedErr bool
	}{
		{
			name: "status good; direct auth",
			req: &v1.PowerRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "10.1.1.1",
							},
							Username: "admin",
							Password: "admin",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				Action:      0,
				SoftTimeout: 0,
				OffDuration: 0,
			},
			message:     "on",
			expectedErr: false,
		},
		{
			name: "status good; external auth",
			req: &v1.PowerRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "10.1.1.1",
							},
							Username: "admin",
							Password: "admin",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				Action:      0,
				SoftTimeout: 0,
				OffDuration: 0,
			},
			message:     "on",
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
			f := freecache.NewStore(freecache.DefaultOptions)
			s := gokv.Store(f)
			var repo repository.Actions
			repo = &persistence.GoKV{
				Store: s,
				Ctx:   ctx,
			}

			var taskRunner task.Task
			taskRunner = &taskrunner.Runner{
				Repository: repo,
				Ctx:        ctx,
				Log:        logger,
			}
			machineSvc := machineService{
				log:        logger,
				taskRunner: taskRunner,
			}
			response, err := machineSvc.powerAction(ctx, testCase.req)

			t.Log("Got response: ", response)
			t.Log("Got err: ", err)

			if testCase.expectedErr {
				g.Expect(response).ToNot(gomega.BeNil(), "Result should be nil")
				g.Expect(err).ToNot(gomega.BeNil(), "Result should be nil")
			} else {
				g.Expect(response.TaskId).Should(gomega.HaveLen(20))
			}
		})
	}
}
