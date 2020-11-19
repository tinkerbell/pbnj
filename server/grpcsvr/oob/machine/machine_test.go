package machine

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/hashicorp/go-multierror"
	"github.com/packethost/pkg/log/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	goipmi "github.com/vmware/goipmi"
)

func TestBootDevice(t *testing.T) {
	sim := goipmi.NewSimulator(net.UDPAddr{})
	err := sim.Run()
	if err != nil {
		t.Fatal(err)
	}
	port := sim.LocalAddr().Port
	defer sim.Stop()
	testCases := []struct {
		name        string
		req         *v1.DeviceRequest
		message     string
		expectedErr error
	}{
		{
			name: "set boot device",
			req: &v1.DeviceRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: fmt.Sprintf("127.0.0.1:%v", port),
							},
							Username: "admin",
							Password: "admin",
						},
					},
				},
				BootDevice: v1.BootDevice_BOOT_DEVICE_BIOS,
			},
			message: "boot device set: bios",
			expectedErr: &multierror.Error{
				Errors: []error{
					errors.New("set boot device failed"),
				},
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			l, zapLogger, _ := logr.NewPacketLogr()
			ctx = ctxzap.ToContext(ctx, zapLogger)
			ma, err := NewMachine(
				WithDeviceRequest(tc.req),
				WithLogger(l),
				WithStatusMessage(make(chan string)),
				WithDeviceRequest(testCase.req),
			)
			if err != nil {
				t.Fatal(err)
			}
			result, errMsg := ma.BootDevice(ctx, testCase.req.BootDevice.String())
			t.Log("result got: ", result)
			t.Log("errMsg got: ", fmt.Sprintf("%+v", errMsg))

			if errMsg != nil {
				diff := cmp.Diff(testCase.expectedErr.Error(), errMsg.Error())
				if diff != "" {
					t.Log(fmt.Sprintf("%+v", errMsg))
					t.Fatalf(diff)
				}
			} else {
				diff := cmp.Diff(testCase.message, result)
				if diff != "" {
					t.Fatalf(diff)
				}
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
