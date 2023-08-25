package taskrunner

import (
	"context"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/tinkerbell/pbnj/pkg/metrics"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

// Runner for executing a task.
type Runner struct {
	Repository repository.Actions
	Ctx        context.Context
	active     atomic.Int32
	total      atomic.Int32
	Dispatcher *dispatcher
}

type dispatcher struct {
	// IngestQueue is a queue of jobs that is process synchronously.
	// It's the entry point for all jobs.
	IngestQueue *IngestQueue
	// perID hold a queue per ID.
	// jobs across different IDs are processed concurrently.
	// jobs with the same ID are processed synchronously.
	perID              sync.Map
	goroutinePerID     atomic.Int32
	idleWorkerShutdown time.Duration
	log                logr.Logger
}

type taskChannel struct {
	ch chan Ingest
}

func NewDispatcher() *dispatcher {
	return &dispatcher{
		IngestQueue:        NewIngestQueue(),
		perID:              sync.Map{},
		idleWorkerShutdown: time.Second * 10,
	}
}

// ActiveWorkers returns a count of currently active worker jobs.
func (r *Runner) ActiveWorkers() int {
	return int(r.active.Load())
}

// TotalWorkers returns a count total workers executed.
func (r *Runner) TotalWorkers() int {
	return int(r.total.Load())
}

func (d *Runner) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			elem, err := d.Dispatcher.IngestQueue.Dequeue()
			if err != nil {
				break
			}

			ch := make(chan Ingest)
			v, loaded := d.Dispatcher.perID.LoadOrStore(elem.Host, taskChannel{ch: ch})
			if !loaded {
				go d.worker1(ctx, elem.Host, ch)
				d.Dispatcher.goroutinePerID.Add(1)
			}
			v.(taskChannel).ch <- elem
		}
	}
}

// channelWorker is a worker that listens on a channel for jobs.
// It will shutdown the worker after gc duration of no elements in the channel or the context is canceled.
// worker is in charge of its own lifecycle.
func (q *Runner) worker1(ctx context.Context, id string, ch chan Ingest) {
	defer func() {
		// do i need to delete the channel if i delete the map entry it lives in?
		// close(ch)
		q.Dispatcher.perID.Delete(id)
		q.Dispatcher.goroutinePerID.Add(-1)
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ch:
			// execute the task synchronously
			q.worker(ctx, q.Dispatcher.log, t.Description, t.ID, t.Action)
		case <-time.After(q.Dispatcher.idleWorkerShutdown):
			// shutdown the worker after this duration of no elements in the channel.
			return
		}
	}
}

// Execute a task, update repository with status.
func (r *Runner) Execute(ctx context.Context, l logr.Logger, description, taskID, host string, action func(chan string) (string, error)) {
	i := Ingest{
		ID:          taskID,
		Host:        host,
		Description: description,
		Action:      action,
	}
	r.Dispatcher.log = l
	r.Dispatcher.IngestQueue.Enqueue(i)
}

func (r *Runner) updateMessages(ctx context.Context, taskID, desc string, ch chan string) error {
	sessionRecord := repository.Record{
		ID:          taskID,
		Description: desc,
		State:       "running",
		Messages:    []string{},
		Error: &repository.Error{
			Code:    0,
			Message: "",
			Details: nil,
		},
	}
	err := r.Repository.Create(taskID, sessionRecord)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-ch:
			record, err := r.Repository.Get(taskID)
			if err != nil {
				return err
			}
			record.Messages = append(record.Messages, msg)
			if err := r.Repository.Update(taskID, record); err != nil {
				return err
			}
		}
	}
}

// does the work, updates the repo record
// TODO handle retrys, use a timeout.
func (r *Runner) worker(ctx context.Context, logger logr.Logger, description, taskID string, action func(chan string) (string, error)) {
	logger = logger.WithValues("taskID", taskID, "description", description)
	r.active.Add(1)
	r.total.Add(1)
	defer func() {
		r.active.Add(-1)
	}()

	metrics.TasksTotal.Inc()
	metrics.TasksActive.Inc()
	defer metrics.TasksActive.Dec()

	messagesChan := make(chan string)
	defer close(messagesChan)
	go r.updateMessages(ctx, taskID, description, messagesChan)

	resultRecord := repository.Record{
		State:    "complete",
		Complete: true,
	}
	result, err := action(messagesChan)
	if err != nil {
		resultRecord.Result = "action failed"
		re, ok := err.(*repository.Error)
		if ok {
			resultRecord.Error = re.StructuredError()
		} else {
			resultRecord.Error.Message = err.Error()
		}
		var foundErr *repository.Error
		if errors.As(err, &foundErr) {
			resultRecord.Error = foundErr.StructuredError()
		}
		logger.Error(err, "task completed with an error")
	}
	record, err := r.Repository.Get(taskID)
	if err != nil {
		return
	}

	record.Complete = resultRecord.Complete
	record.State = resultRecord.State
	record.Result = result
	record.Error = resultRecord.Error

	if err := r.Repository.Update(taskID, record); err != nil {
		logger.Error(err, "failed to update record")
	}

	logger.Info("worker complete", "complete", true)
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
