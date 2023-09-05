package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/gomega"
	v1 "github.com/tinkerbell/pbnj/api/v1"
)

func TestTaskFound(t *testing.T) {
	pr := &v1.PowerRequest{
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
		PowerAction: v1.PowerAction_POWER_ACTION_STATUS,
		SoftTimeout: 0,
		OffDuration: 0,
	}
	resp, err := machineService.Power(context.Background(), pr)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	time.Sleep(time.Second)
	taskReq := &v1.StatusRequest{TaskId: resp.TaskId}
	taskResp, _ := taskService.Status(context.Background(), taskReq)
	if taskResp.Id != resp.TaskId {
		t.Fatalf("got: %+v", taskResp)
	}
}

func TestRecordNotFound(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.StatusRequest
		message     string
		expectedErr bool
	}{
		{
			name:        "record of task not found",
			req:         &v1.StatusRequest{TaskId: "123"},
			message:     "rpc error: code = NotFound desc = record id not found: 123",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)

			ctx := context.Background()
			response, err := taskService.Status(ctx, testCase.req)

			t.Log("Got response: ", response)
			t.Log("Got err: ", err)

			if testCase.expectedErr {
				g.Expect(response).To(gomega.BeNil(), "Response should be nil")
				g.Expect(err).ToNot(gomega.BeNil(), "error should not be nil")
				g.Expect(err.Error()).To(gomega.Equal(testCase.message))
			} else {
				g.Expect(response.Result).To(gomega.Equal(testCase.message))
			}
		})
	}
}
