package rpc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/persistence"
	"github.com/tinkerbell/pbnj/grpc/taskrunner"
)

const tempIPMITool = "/tmp/ipmitool"

var (
	tr             *taskrunner.Runner
	bmcService     BmcService
	taskService    TaskService
	machineService MachineService
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	repo := &persistence.GoKV{
		Store: s,
		Ctx:   ctx,
	}

	tr = taskrunner.NewRunner(repo, 100, time.Second)
	tr.Start(ctx)
	bmcService = BmcService{
		TaskRunner:             tr,
		UnimplementedBMCServer: v1.UnimplementedBMCServer{},
	}
	taskService = TaskService{
		TaskRunner: tr,
	}
	machineService = MachineService{
		TaskRunner: tr,
	}
	_, err := exec.LookPath("ipmitool")
	if err != nil {
		err := os.WriteFile(tempIPMITool, []byte{}, 0o777)
		if err != nil {
			fmt.Println("didnt find ipmitool in PATH and couldnt create one in /tmp")
			os.Exit(3) //nolint:revive // deep-exit here is OK
		}
		path := os.Getenv("PATH")
		os.Setenv("PATH", fmt.Sprintf("%v:/tmp", path))
	}
}

func teardown() {
	os.Remove(tempIPMITool)
}

func TestConfigNetworkSource(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.NetworkSourceRequest
		message     string
		expectedErr error
	}{
		{
			name: "status good",
			req: &v1.NetworkSourceRequest{
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
				NetworkSource: 0,
			},
			message:     "good",
			expectedErr: errors.New("not implemented"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			response, err := bmcService.NetworkSource(context.Background(), testCase.req)
			if response != nil {
				t.Fatalf("response should be nil, got: %v", response)
			}
			if diff := cmp.Diff(tc.expectedErr.Error(), err.Error()); diff != "" {
				t.Fatal(diff)
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
			response, err := bmcService.Reset(context.Background(), tc.in)
			if err != nil {
				diff := cmp.Diff(tc.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else if response.TaskId == "" {
				t.Fatal("expected taskId, got:", response.TaskId)
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
			response, err := bmcService.CreateUser(context.Background(), tc.in)
			if err != nil {
				diff := cmp.Diff(tc.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else if response.TaskId == "" {
				t.Fatal("expected taskId, got:", response.TaskId)
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
			response, err := bmcService.UpdateUser(context.Background(), tc.in)
			if err != nil {
				diff := cmp.Diff(tc.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else if response.TaskId == "" {
				t.Fatal("expected taskId, got:", response.TaskId)
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
			response, err := bmcService.DeleteUser(context.Background(), tc.in)
			if err != nil {
				diff := cmp.Diff(tc.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else if response.TaskId == "" {
				t.Fatal("expected taskId, got:", response.TaskId)
			}
		})
	}
}
