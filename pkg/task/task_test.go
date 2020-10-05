package task

import (
	"context"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging/zaplog"
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

func TestRoundTrip(t *testing.T) {
	description := "test task"
	defaultError := &oob.Error{
		Error: &v1.Error{
			Code:    0,
			Message: "",
			Details: nil,
		},
	}
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	defer s.Close()
	repo := &repository.GoKV{Store: s}
	runner := Runner{Repository: repo}

	logger, zapLogger, _ := zaplog.RegisterLogger()
	ctx = ctxzap.ToContext(ctx, zapLogger)
	id, err := runner.Execute(ctx, logger, description, func(s chan string) (string, *oob.Error) {
		return "didnt do anything", defaultError
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 20 {
		t.Fatalf("expected id of length 20,  got: %v (%v)", len(id), id)
	}

	// must be min of 3 because we sleep 2 seconds in worker function to allow final status messages to be written
	time.Sleep(500 * time.Millisecond)
	record, err := runner.Status(ctx, logger, id)
	if err != nil {
		t.Fatal(err)
	}
	if record.StatusResponse.Complete != true {
		t.Fatalf("expected task to be complete, got: %+v", record)
	}
}
