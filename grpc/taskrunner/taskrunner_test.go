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
	runner := NewRunner(repo, 100, time.Second)
	runner.Start(ctx)
	time.Sleep(time.Millisecond * 100)

	taskID := xid.New().String()
	if len(taskID) != 20 {
		t.Fatalf("expected id of length 20,  got: %v (%v)", len(taskID), taskID)
	}
	runner.Execute(ctx, logger, description, taskID, "123", func(s chan string) (string, error) {
		return "didnt do anything", defaultError
	})

	// must be min of 3 because we sleep 2 seconds in worker function to allow final status messages to be written
	time.Sleep(time.Second * 2)
	record, err := runner.Status(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}

	if !record.Complete {
		t.Fatalf("expected task to be complete, got: %+v", record)
	}
}
