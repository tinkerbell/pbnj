package taskrunner

import (
	"context"
	"net"
	"net/url"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// Runner for executing a task.
type Runner struct {
	Repository repository.Actions
	Ctx        context.Context
	Log        logging.Logger
	active     int
	total      int
	counterMu  sync.RWMutex
}

// ActiveWorkers returns a count of currently active worker jobs.
func (r *Runner) ActiveWorkers() int {
	r.counterMu.RLock()
	defer r.counterMu.RUnlock()
	return r.active
}

// TotalWorkers returns a count total workers executed.
func (r *Runner) TotalWorkers() int {
	r.counterMu.RLock()
	defer r.counterMu.RUnlock()
	return r.total
}

// Execute a task, update repository with status.
func (r *Runner) Execute(ctx context.Context, description, taskID string, action func(chan string) (string, error)) {
	go r.worker(ctx, description, taskID, action)
}

// does the work, updates the repo record
// TODO handle retrys, use a timeout.
func (r *Runner) worker(ctx context.Context, description, taskID string, action func(chan string) (string, error)) {
	logger := r.Log.GetContextLogger(ctx)
	logger = logger.WithValues("complete", false, "taskID", taskID, "description", description)
	r.counterMu.Lock()
	r.active++
	r.total++
	r.counterMu.Unlock()
	defer func() {
		r.counterMu.Lock()
		r.active--
		r.counterMu.Unlock()
	}()

	metrics.TasksTotal.Inc()
	metrics.TasksActive.Inc()
	defer metrics.TasksActive.Dec()

	messagesChan := make(chan string)
	actionACK := make(chan bool, 1)
	actionSyn := make(chan bool, 1)
	defer close(messagesChan)
	defer close(actionACK)
	defer close(actionSyn)
	repo := r.Repository
	sessionRecord := repository.Record{
		ID:          taskID,
		Description: description,
		State:       "running",
		Messages:    []string{},
		Error: &repository.Error{
			Code:    0,
			Message: "",
			Details: nil,
		},
	}

	err := repo.Create(taskID, sessionRecord)
	if err != nil {
		// TODO how to handle unable to create record; ie network error, persistence error, etc?
		logger.V(0).Error(err, "task complete", "complete", true)
		return
	}

	go func() {
		for {
			select {
			case msg := <-messagesChan:
				currStatus, _ := repo.Get(taskID)
				sessionRecord.Messages = append(currStatus.Messages, msg) // nolint:gocritic // apparently this is the right slice
				_ = repo.Update(taskID, sessionRecord)
			case <-actionSyn:
				actionACK <- true
				return
			default:
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	sessionRecord.Result, err = action(messagesChan)
	actionSyn <- true
	<-actionACK
	sessionRecord.State = "complete"
	sessionRecord.Complete = true
	var finalErr error
	if err != nil {
		finalErr = multierror.Append(finalErr, err)
		sessionRecord.Result = "action failed"
		re, ok := err.(*repository.Error)
		if ok {
			sessionRecord.Error = re.StructuredError()
		} else {
			sessionRecord.Error.Message = err.Error()
		}
		var foundErr *repository.Error
		if errors.As(err, &foundErr) {
			sessionRecord.Error = foundErr.StructuredError()
		}
	}
	// TODO handle unable to update record; ie network error, persistence error, etc
	if err := repo.Update(taskID, sessionRecord); err != nil {
		finalErr = multierror.Append(finalErr, err)
	}

	if finalErr != nil {
		logger.Error(finalErr, "task complete", "complete", true)
	} else {
		logger.Info("task complete", "complete", true)
	}
}

// Status returns the status record of a task.
func (r *Runner) Status(_ context.Context, taskID string) (record repository.Record, err error) {
	record, err = r.Repository.Get(taskID)
	if err != nil {
		switch t := err.(type) {
		case *net.OpError:
			if t.Op == "dial" {
				return record, errors.New("persistence error: unknown host")
			} else if t.Op == "read" {
				return record, errors.New("persistence error: connection refused")
			}
		case syscall.Errno:
			if t == syscall.ECONNREFUSED {
				return record, errors.New("persistence error: connection refused")
			}
		case *url.Error:
			return record, errors.New("persistence error: connection refused")
		}
	}
	return
}
