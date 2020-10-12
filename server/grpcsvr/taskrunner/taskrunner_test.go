package taskrunner

import (
	"context"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	"github.com/tinkerbell/pbnj/cmd/zaplog"
	v1 "github.com/tinkerbell/pbnj/pkg/api/v1"
	"github.com/tinkerbell/pbnj/pkg/oob"
	"github.com/tinkerbell/pbnj/pkg/repository"
	"github.com/tinkerbell/pbnj/server/grpcsvr/persistence"
)

func TestRoundTrip(t *testing.T) {
	description := "test task"
	defaultError := oob.Error{
		Error: v1.Error{
			Code:    0,
			Message: "",
			Details: nil,
		},
	}
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	defer s.Close()
	var repo repository.Actions
	repo = &persistence.GoKV{Store: s, Ctx: ctx}
	logger, zapLogger, _ := zaplog.RegisterLogger()
	ctx = ctxzap.ToContext(ctx, zapLogger)
	runner := Runner{
		Repository: repo,
		Ctx:        ctx,
		Log:        logger,
	}

	id, err := runner.Execute(description, func(s chan string) (string, oob.Error) {
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
	record, err := runner.Status(id)
	if err != nil {
		t.Fatal(err)
	}
	if record.StatusResponse.Complete != true {
		t.Fatalf("expected task to be complete, got: %+v", record)
	}
}
