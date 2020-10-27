package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/onsi/gomega"
	"github.com/packethost/pkg/log/logr"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/cmd/zaplog"
	"github.com/tinkerbell/pbnj/pkg/repository"
	"github.com/tinkerbell/pbnj/server/grpcsvr/persistence"
	"github.com/tinkerbell/pbnj/server/grpcsvr/taskrunner"
)

func TestTaskFound(t *testing.T) {
	// create a task
	ctx := context.Background()
	defaultError := repository.Error{
		Code:    0,
		Message: "",
		Details: nil,
	}
	l, zapLogger, _ := logr.NewPacketLogr()
	logger := zaplog.RegisterLogger(l)
	ctx = ctxzap.ToContext(ctx, zapLogger)
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	repo := &persistence.GoKV{Store: s, Ctx: ctx}

	taskRunner := &taskrunner.Runner{
		Repository: repo,
		Ctx:        ctx,
		Log:        logger,
	}
	taskID, err := taskRunner.Execute("test", func(s chan string) (string, repository.Error) {
		return "doing cool stuff", defaultError // nolint
	})
	if err != nil {
		t.Fatal(err)
	}
	taskReq := &v1.StatusRequest{TaskId: taskID}

	taskSvc := TaskService{
		Log:        logger,
		TaskRunner: taskRunner,
	}

	time.Sleep(10 * time.Millisecond)
	taskResp, err := taskSvc.Task(ctx, taskReq)
	if err != nil {
		t.Fatal(err)
	}
	if taskResp.Id != taskID {
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
			message:     "record id not found: 123",
			expectedErr: true,
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
			repo := &persistence.GoKV{Store: s, Ctx: ctx}

			taskRunner := &taskrunner.Runner{
				Repository: repo,
				Ctx:        ctx,
				Log:        logger,
			}
			taskSvc := TaskService{
				Log:        logger,
				TaskRunner: taskRunner,
			}
			response, err := taskSvc.Task(ctx, testCase.req)

			t.Log("Got : ", response)
			t.Log("Got : ", err)

			if testCase.expectedErr {
				g.Expect(response).To(gomega.BeNil(), "Result should be nil")
				g.Expect(err).ToNot(gomega.BeNil(), "Result should be nil")
				g.Expect(err.Error()).To(gomega.Equal(testCase.message))
			} else {
				g.Expect(response.Result).To(gomega.Equal(testCase.message))
			}
		})
	}
}
