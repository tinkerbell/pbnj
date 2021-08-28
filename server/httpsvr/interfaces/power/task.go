// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package power

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/reqid"
)

// TaskTimeout is the default timeout for tasks.
const TaskTimeout = 3 * time.Minute

var tasks struct {
	sync.Mutex
	m map[string]*Task
}

// FindTask fetches the Taskf for the provided id.
func FindTask(id string) *Task {
	tasks.Lock()
	defer tasks.Unlock()

	if tasks.m == nil {
		return nil
	}
	return tasks.m[id]
}

// DeleteTask deletes the Task with the corresponding id.
func DeleteTask(id string) {
	tasks.Lock()
	defer tasks.Unlock()

	if tasks.m == nil {
		return
	}
	delete(tasks.m, id)
}

// CleanupTasks deletes expired tasks.
func CleanupTasks(expiration time.Duration) {
	var (
		ids   []string
		tx    = elog.TxFromContext(reqid.WithID(context.Background(), "cleanup"))
		start = time.Now()
	)
	defer func() {
		duration := time.Since(start)
		if len(ids) > 0 {
			tx.Notice("cleanup", "tasks", strings.Join(ids, ","), "duration", duration)
		} else {
			tx.Debug("cleanup", "duration", duration)
		}
	}()

	tasks.Lock()
	defer tasks.Unlock()

	for _, t := range tasks.m {
		if start.Sub(t.start) > expiration {
			delete(tasks.m, t.id)
			ids = append(ids, t.id)
		}
	}
}

// StartTask starts a new task.
func StartTask(ctx context.Context, id string, op Operation, driver Driver, opts Options) *Task {
	t := newTask(id, op)
	if t == nil {
		return nil
	}
	t.start = time.Now()
	go t.run(ctx, driver, opts)
	return t
}

// Task represents a task.
type Task struct {
	id     string
	op     Operation
	start  time.Time
	end    time.Time
	done   chan struct{}
	status Status
	err    error
}

func newTask(id string, op Operation) *Task {
	t := &Task{
		id:   id,
		op:   op,
		done: make(chan struct{}),
	}

	tasks.Lock()
	defer tasks.Unlock()

	if tasks.m == nil {
		tasks.m = make(map[string]*Task)
	} else if _, ok := tasks.m[id]; ok {
		return nil
	}
	tasks.m[id] = t
	return t
}

// Done returns a channel that will remain blocked while the task is in progress.
func (t *Task) Done() <-chan struct{} {
	return t.done
}

// Err returns the first non-nil error encountered by Task.
func (t *Task) Err() error {
	return errors.WithMessage(t.err, "task error")
}

// ID returns the task id.
func (t *Task) ID() string {
	return t.id
}

func (t *Task) run(ctx context.Context, driver Driver, opts Options) {
	ctx, cancel := context.WithTimeout(ctx, TaskTimeout)
	defer cancel()

	// cleanup
	defer func() {
		t.end = time.Now()

		if iface := recover(); iface != nil {
			var err error
			switch iface := iface.(type) {
			case error:
				err = iface
			default:
				err = errors.Errorf("%v", iface)
			}
			err = errors.WithMessage(err, "panic recovered")

			t.err = err
		}

		t.status = driver.LastStatus()
		err := driver.Close()
		if !opts.IgnoreRunError && t.err != nil {
			t.err = err
		}

		close(t.done)
	}()

	err := t.op(ctx, driver, opts)
	if !opts.IgnoreRunError {
		t.err = err
	}
}
