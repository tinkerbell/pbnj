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
	Repository   repository.Actions
	active       atomic.Int32
	total        atomic.Int32
	orchestrator *orchestrator
}

type Task struct {
	ID          string                            `json:"id"`
	Host        string                            `json:"host"`
	Description string                            `json:"description"`
	Action      func(chan string) (string, error) `json:"-"`
	Log         logr.Logger                       `json:"-"`
}

// NewRunner returns a task runner that manages tasks, workers, queues, and persistence.
//
// maxIngestionWorkers is the maximum number of concurrent workers that will be allowed.
// These are the workers that handle ingesting tasks from RPC endpoints and writing them to the map of per Host ID queues.
//
// maxWorkers is the maximum number of concurrent workers that will be allowed to handle bmc tasks.
//
// workerIdleTimeout is the idle timeout for workers. If no tasks are received within the timeout, the worker will exit.
func NewRunner(repo repository.Actions, maxWorkers int, workerIdleTimeout time.Duration) *Runner {
	o := &orchestrator{
		workers:  sync.Map{},
		fifoChan: make(chan string, 10000),
		// perHostChan is a map of hostID to a channel of tasks.
		perHostChan:       sync.Map{},
		manager:           newManager(maxWorkers),
		workerIdleTimeout: workerIdleTimeout,
		ingestChan:        make(chan Task, 10000),
	}

	return &Runner{
		Repository:   repo,
		orchestrator: o,
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

func (r *Runner) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metrics.NumWorkers.Set(float64(r.orchestrator.manager.RunningCount()))
				var size int
				r.orchestrator.workers.Range(func(key, value interface{}) bool {
					size++
					return true
				})
				metrics.WorkerMap.Set(float64(size))
			}
		}
	}()
	go r.ingest(ctx)
	go r.orchestrate(ctx)
}

// Execute a task, update repository with status.
func (r *Runner) Execute(_ context.Context, l logr.Logger, description, taskID, hostID string, action func(chan string) (string, error)) {
	i := Task{
		ID:          taskID,
		Host:        hostID,
		Description: description,
		Action:      action,
		Log:         l,
	}

	r.orchestrator.ingestChan <- i
	metrics.IngestionQueue.Inc()
	metrics.Ingested.Inc()
}

func (r *Runner) updateMessages(ctx context.Context, taskID string, ch chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			record, err := r.Repository.Get(taskID)
			if err != nil {
				return
			}
			record.Messages = append(record.Messages, msg)
			if err := r.Repository.Update(taskID, record); err != nil {
				return
			}
		}
	}
}

// does the work, updates the repo record
// TODO handle retrys, use a timeout.
func (r *Runner) process(ctx context.Context, logger logr.Logger, description, taskID string, action func(chan string) (string, error)) {
	logger = logger.WithValues("taskID", taskID, "description", description)
	r.active.Add(1)
	r.total.Add(1)
	defer func() {
		r.active.Add(-1)
	}()

	defer metrics.TasksTotal.Inc()
	defer metrics.TotalGauge.Inc()
	metrics.TasksActive.Inc()
	defer metrics.TasksActive.Dec()

	messagesChan := make(chan string)
	defer close(messagesChan)
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
	err := r.Repository.Create(taskID, sessionRecord)
	if err != nil {
		return
	}
	cctx, done := context.WithCancel(ctx)
	defer done()
	go r.updateMessages(cctx, taskID, messagesChan)

	resultRecord := repository.Record{
		State:    "complete",
		Complete: true,
		Error: &repository.Error{
			Code:    0,
			Message: "",
			Details: nil,
		},
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
