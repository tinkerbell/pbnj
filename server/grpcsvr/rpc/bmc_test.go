package rpc

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/onsi/gomega"
	"github.com/tinkerbell/pbnj/cmd/zaplog"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
)

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
			t.Parallel()
			g := gomega.NewGomegaWithT(t)

			ctx := context.Background()

			logger, zapLogger, _ := zaplog.RegisterLogger()
			ctx = ctxzap.ToContext(ctx, zapLogger)
			bmcSvc := BmcService{
				Log: logger,
			}
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