package machine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/packethost/pkg/log/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

func TestBootDevice(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.DeviceRequest
		message     string
		expectedErr repository.Error
	}{
		{
			name: "timeout",
			req: &v1.DeviceRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "127.0.0.1",
							},
							Username: "admin",
							Password: "admin",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "none",
				},
				Persistent: false,
				EfiBoot:    false,
			},
			message: "good",
			expectedErr: repository.Error{
				Code:    2,
				Message: "could not connect",
				Details: []string{"context deadline exceeded"},
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
			defer cancel()

			l, zapLogger, _ := logr.NewPacketLogr()
			ctx = ctxzap.ToContext(ctx, zapLogger)

			ma, err := NewMachine(
				WithDeviceRequest(tc.req),
				WithContext(ctx),
				WithLogger(l),
				WithStatusMessage(make(chan string)),
			)
			if err != nil {
				t.Fatal(err)
			}
			result, errMsg := ma.BootDevice(ctx)
			t.Log("result got: ", result)
			t.Log("errMsg got: ", fmt.Sprintf("%+v", errMsg))

			diff := cmp.Diff(testCase.expectedErr, errMsg)
			if diff != "" {
				t.Log(fmt.Sprintf("%+v", errMsg))
				t.Fatalf(diff)
			}
		})
	}
}

/*
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

			l, zapLogger, _ := logr.NewPacketLogr()
			logger := zaplog.RegisterLogger(l)
			ctx = ctxzap.ToContext(ctx, zapLogger)
			f := freecache.NewStore(freecache.DefaultOptions)
			s := gokv.Store(f)
			repo := &persistence.GoKV{
				Store: s,
				Ctx:   ctx,
			}

			taskRunner := &taskrunner.Runner{
				Repository: repo,
				Ctx:        ctx,
				Log:        logger,
			}
			machineSvc := MachineService{
				Log:        logger,
				TaskRunner: taskRunner,
			}
			response, err := machineSvc.PowerAction(ctx, testCase.req)

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
*/
