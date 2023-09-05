package rpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/persistence"
	"github.com/tinkerbell/pbnj/grpc/taskrunner"
)

func TestDevice(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.DeviceRequest
		message     string
		expectedErr error
	}{
		{
			name: "status good; direct auth",
			req: &v1.DeviceRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "127.0.0.1",
							},
							Username: "ADMIN",
							Password: "ADMIN",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
				Persistent: false,
				EfiBoot:    false,
			},
			message: "good",
		},
		{
			name:        "validation failure",
			req:         &v1.DeviceRequest{Authn: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{}}}},
			message:     "",
			expectedErr: errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Username: value '' must not be an empty string"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)

			ctx := context.Background()

			f := freecache.NewStore(freecache.DefaultOptions)
			s := gokv.Store(f)
			repo := &persistence.GoKV{
				Store: s,
				Ctx:   ctx,
			}

			tr := taskrunner.NewRunner(repo)
			go tr.Start(ctx)
			machineSvc := MachineService{
				TaskRunner: tr,
			}
			time.Sleep(time.Millisecond * 100)
			response, err := machineSvc.BootDevice(ctx, testCase.req)

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

func TestPower(t *testing.T) {
	// TODO make sure external auth doesnt break stuff
	testCases := []struct {
		name        string
		req         *v1.PowerRequest
		message     string
		expectedErr error
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
				PowerAction: 0,
				SoftTimeout: 0,
				OffDuration: 0,
			},
			message: "on",
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
				PowerAction: 0,
				SoftTimeout: 0,
				OffDuration: 0,
			},
			message: "on",
		},
		{
			name:        "validation failure",
			req:         &v1.PowerRequest{Authn: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{}}}},
			message:     "",
			expectedErr: errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Username: value '' must not be an empty string"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)

			ctx := context.Background()

			f := freecache.NewStore(freecache.DefaultOptions)
			s := gokv.Store(f)
			repo := &persistence.GoKV{
				Store: s,
				Ctx:   ctx,
			}

			tr := taskrunner.NewRunner(repo)
			go tr.Start(ctx)
			time.Sleep(time.Millisecond * 100)
			machineSvc := MachineService{
				TaskRunner: tr,
			}
			response, err := machineSvc.Power(ctx, testCase.req)

			t.Log("Got response: ", response)
			t.Log("Got err: ", err)
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
