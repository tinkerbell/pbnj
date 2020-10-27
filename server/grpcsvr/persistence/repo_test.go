package persistence

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

func TestAllMethods(t *testing.T) {
	id := "1234567"
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	defer s.Close()
	repo := &GoKV{Store: s, Ctx: ctx}
	record := repository.Record{
		Id:          id,
		Description: "test record",
		Error:       nil,
		State:       "running",
		Result:      "",
		Complete:    false,
	}
	expected := repository.Record{
		Id:          id,
		Description: "test record",
		Error:       nil,
		State:       "running",
		Result:      "",
		Complete:    false,
	}
	expectedUpdated := repository.Record{
		Id:          id,
		Description: "test record",
		Error:       nil,
		State:       "complete",
		Result:      "did a good thing",
		Complete:    true,
	}
	err := repo.Create(id, record)
	if err != nil {
		t.Fatal(err)
	}

	result, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(expected, result); diff != nil {
		t.Fatal(diff)
	}

	// update record
	result.Complete = true
	result.Result = "did a good thing"
	result.State = "complete"
	err = repo.Update(id, result)
	if err != nil {
		t.Fatal(err)
	}

	updatedResult, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(expectedUpdated, updatedResult); diff != nil {
		t.Fatal(diff)
	}

	// delete record
	err = repo.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetRecordNotFound(t *testing.T) {
	id := "123"
	expectedError := fmt.Sprintf("record id not found: %v", id)
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	repo := &GoKV{Store: s, Ctx: ctx}
	_, err := repo.Get(id)
	if err == nil {
		t.Fatalf("expecting NON nil error")
	}
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected: %v, got: %v", expectedError, err.Error())
	}
}

func TestGetRecordError(t *testing.T) {
	id := ""
	expectedError := "The passed key is an empty string, which is invalid"
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	repo := &GoKV{Store: s, Ctx: ctx}
	_, err := repo.Get(id)
	if err == nil {
		t.Fatalf("expecting NON nil error")
	}
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected: %v, got: %v", expectedError, err.Error())
	}
}

func TestUpdateRecordNotFound(t *testing.T) {
	id := "123"
	expectedError := fmt.Sprintf("record id not found: %v", id)
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	repo := &GoKV{Store: s, Ctx: ctx}
	err := repo.Update(id, repository.Record{})
	if err == nil {
		t.Fatalf("expecting NON nil error")
	}
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected: %v, got: %v", expectedError, err.Error())
	}
}

func TestUpdateRecordError(t *testing.T) {
	id := ""
	expectedError := "The passed key is an empty string, which is invalid"
	ctx := context.Background()
	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	repo := &GoKV{Store: s, Ctx: ctx}
	err := repo.Update(id, repository.Record{})
	if err == nil {
		t.Fatalf("expecting NON nil error")
	}
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected: %v, got: %v", expectedError, err.Error())
	}
}
