package machine

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/bmc-toolbox/bmclib"
	"github.com/google/go-cmp/cmp"
	"github.com/jacobweinstock/registrar"
	"github.com/packethost/pkg/log/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
	common "github.com/tinkerbell/pbnj/server/grpcsvr/oob"
)

func newAction(withAuthErr bool) Action {
	packetLogr, _, _ := logr.NewPacketLogr()
	var authn *v1.Authn_DirectAuthn
	if withAuthErr {
		authn = &v1.Authn_DirectAuthn{
			DirectAuthn: nil,
		}
	} else {
		authn = &v1.Authn_DirectAuthn{
			DirectAuthn: &v1.DirectAuthn{
				Host: &v1.Host{
					Host: "localhost",
				},
				Username: "admin",
				Password: "admin",
			},
		}
	}
	m := Action{
		Accessory: common.Accessory{
			Log:            packetLogr.Logger,
			StatusMessages: make(chan string),
		},
		BootDeviceRequest: &v1.DeviceRequest{
			Authn: &v1.Authn{
				Authn: authn,
			},
			Vendor: &v1.Vendor{
				Name: "local",
			},
			BootDevice: v1.BootDevice_BOOT_DEVICE_PXE,
		},
	}
	return m
}

func TestBootDevice(t *testing.T) {
	b := bmclib.NewClient("localhost", "623", "admin", "admin")
	m := newAction(false)
	authErr := newAction(true)

	testCases := []struct {
		name         string
		ok           bool
		err          error
		want         string
		wantErr      error
		bootDevice   string
		actionStruct Action
	}{
		{"reset err", false, errors.New("bad"), "", &repository.Error{Code: v1.Code_value["UNKNOWN"], Message: "bad", Details: []string{}}, v1.BootDevice_BOOT_DEVICE_PXE.String(), m},
		{"success", true, nil, "", nil, v1.BootDevice_BOOT_DEVICE_PXE.String(), m},
		{"reset not ok", false, nil, "", &repository.Error{Code: v1.Code_value["UNKNOWN"], Message: "setting boot device failed", Details: []string{}}, v1.BootDevice_BOOT_DEVICE_PXE.String(), m},
		{"unknown reset request", true, nil, "", &repository.Error{Code: v1.Code_value["INVALID_ARGUMENT"], Message: "unknown boot device", Details: []string{}}, "blah", m},
		{"auth parse err", true, nil, "", &repository.Error{Code: v1.Code_value["UNAUTHENTICATED"], Message: "no auth found", Details: []string{}}, v1.BootDevice_BOOT_DEVICE_PXE.String(), authErr},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(b), "Open", func(_ *bmclib.Client, _ context.Context) (err error) {
				return nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&registrar.Registry{}), "FilterForCompatible", func(_ *registrar.Registry, _ context.Context) (drvs registrar.Drivers) {
				return b.Registry.Drivers
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(b), "SetBootDevice", func(_ *bmclib.Client, _ context.Context, _ string, _ bool, _ bool) (ok bool, err error) {
				return tc.ok, tc.err
			})
			result, err := tc.actionStruct.BootDeviceSet(context.Background(), tc.bootDevice, false, false)
			if err != nil {
				if tc.wantErr != nil {
					diff := cmp.Diff(err.Error(), tc.wantErr.Error())
					if diff != "" {
						t.Fatal(diff)
					}
				}
			} else {
				diff := cmp.Diff(result, tc.want)
				if diff != "" {
					t.Fatal(diff)
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
