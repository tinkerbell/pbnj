package taskrunner

import (
	"context"
	"net"
	"net/url"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/rs/xid"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// Runner for executing a task
type Runner struct {
	Repository repository.Actions
	Ctx        context.Context
	Log        logging.Logger
	active     int
	total      int
	counterMu  sync.RWMutex
}

// ActiveWorkers returns a count of currently active worker jobs
func (r *Runner) ActiveWorkers() int {
	r.counterMu.RLock()
	defer r.counterMu.RUnlock()
	return r.active
}

// TotalWorkers returns a count total workers executed
func (r *Runner) TotalWorkers() int {
	r.counterMu.RLock()
	defer r.counterMu.RUnlock()
	return r.total
}

// Execute a task, update repository with status
func (r *Runner) Execute(ctx context.Context, description string, action func(chan string) (string, error)) (id string, err error) {
	rawID := xid.New()
	id = rawID.String()
	l := r.Log.GetContextLogger(ctx)
	l.V(0).Info("executing task", "taskID", id, "taskDescription", description)
	go r.worker(ctx, l, id, description, action)
	return id, err
}

// does the work, updates the repo record
// TODO handle retrys, use a timeout
func (r *Runner) worker(ctx context.Context, logger logr.Logger, id string, description string, action func(chan string) (string, error)) {
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
	logger.V(0).Info("starting worker", "taskID", id, "description", description)

	messagesChan := make(chan string)
	actionACK := make(chan bool, 1)
	actionSyn := make(chan bool, 1)
	defer close(messagesChan)
	defer close(actionACK)
	defer close(actionSyn)
	repo := r.Repository
	sessionRecord := repository.Record{
		Id:          id,
		Description: description,
		State:       "running",
		Messages:    []string{},
		Error: &repository.Error{
			Code:    0,
			Message: "",
			Details: nil,
		}}

	err := repo.Create(id, sessionRecord)
	if err != nil {
		// TODO how to handle unable to create record; ie network error, persistence error, etc?
		logger.V(0).Error(err, "creating record failed")
		return
	}

	go func() {
		for {
			select {
			case msg := <-messagesChan:
				logger.V(0).Info("STATUS MESSAGE", "statusMsg", msg)
				currStatus, _ := repo.Get(id)
				sessionRecord.Messages = append(currStatus.Messages, msg)
				_ = repo.Update(id, sessionRecord)
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
	if err != nil {
		logger.V(0).Info("error running action", "err", err.Error())
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
	errI := repo.Update(id, sessionRecord)
	if errI != nil {
		// TODO handle unable to update record; ie network error, persistence error, etc
		logger.V(0).Error(err, "updating record failed")
		return
	}
}

// Status returns the status record of a task
func (r *Runner) Status(ctx context.Context, id string) (record repository.Record, err error) {
	l := r.Log.GetContextLogger(ctx)
	l.V(0).Info("getting task record", "taskID", id)
	record, err = r.Repository.Get(id)
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
