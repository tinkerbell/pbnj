package taskrunner

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	"github.com/rs/xid"
	"github.com/tinkerbell/pbnj/grpc/persistence"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

func TestRoundTrip(t *testing.T) {
	description := "test task"
	defaultError := &repository.Error{
		Code:    0,
		Message: "",
		Details: nil,
	}
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	defer s.Close()
	repo := &persistence.GoKV{Store: s, Ctx: ctx}
	logger := logr.Discard()
	runner := Runner{
		Repository: repo,
		Ctx:        ctx,
	}

	taskID := xid.New().String()
	runner.Execute(ctx, logger, description, taskID, func(_ chan string) (string, error) {
		return "didnt do anything", defaultError
	})

	if len(taskID) != 20 {
		t.Fatalf("expected id of length 20,  got: %v (%v)", len(taskID), taskID)
	}

	// must be min of 3 because we sleep 2 seconds in worker function to allow final status messages to be written
	time.Sleep(500 * time.Millisecond)
	record, err := runner.Status(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}

	if !record.Complete {
		t.Fatalf("expected task to be complete, got: %+v", record)
	}
}
