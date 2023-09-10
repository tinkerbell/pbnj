package taskrunner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tinkerbell/pbnj/pkg/metrics"
)

type orchestrator struct {
	workers           sync.Map
	manager           *concurrencyManager
	workerIdleTimeout time.Duration
	fifoChan          chan string
	// perHostChan is a map of hostID to a channel of tasks.
	perHostChan sync.Map
	ingestChan  chan Task
}

// ingest take a task off the ingestion queue and puts it on the perID queue
// and adds the host ID to the fcfs queue.
func (r *Runner) ingest(ctx context.Context) {
	// dequeue from ingestion queue
	// enqueue to perID queue
	// enqueue to fcfs queue
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-r.orchestrator.ingestChan:

			// 2. enqueue to perID queue
			ch := make(chan Task, 10)
			q, exists := r.orchestrator.perHostChan.LoadOrStore(t.Host, ch)
			v, ok := q.(chan Task)
			if !ok {
				fmt.Println("bad type: IngestQueue")
				return
			}
			if exists {
				close(ch)
			}
			v <- t
			metrics.PerIDQueue.WithLabelValues(t.Host).Inc()
			metrics.IngestionQueue.Dec()
			metrics.NumPerIDEnqueued.Inc()
			r.orchestrator.workers.Store(t.Host, false)
		}
	}
}

// 1. dequeue from fcfs queue
// 2. dequeue from perID queue
// 3. if worker id exists, send task to worker. else continue.
// 4. if maxWorkers is reached, wait for available worker. else create worker and send task to worker.
func (r *Runner) orchestrate(ctx context.Context) {
	// 1. dequeue from fcfs queue
	// 2. start workers
	for {
		// time.Sleep(time.Second * 3) - this potential helps with ingestion
		r.orchestrator.workers.Range(func(key, value interface{}) bool {
			// if worker id exists in o.workers, then move on because the worker is already running.
			if value.(bool) { //nolint: forcetypeassert // values are always certain.
				return true
			}

			// wait for a worker to become available
			r.orchestrator.manager.Wait()

			r.orchestrator.workers.Store(key.(string), true) //nolint: forcetypeassert // values are always certain.
			v, found := r.orchestrator.perHostChan.Load(key.(string))
			if !found {
				return false
			}
			go r.worker(ctx, key.(string), v.(chan Task)) //nolint: forcetypeassert // values are always certain.
			return true
		})
	}
}

func (r *Runner) worker(ctx context.Context, hostID string, q chan Task) {
	defer r.orchestrator.manager.Done()
	defer r.orchestrator.workers.Delete(hostID)

	for {
		select {
		case <-ctx.Done():
			// TODO: check queue length before returning maybe?
			// For 175000 tasks, i found there would occasionally be 1 or 2 that didnt get processed.
			// still seemed to be in the queue/chan.
			return
		case t := <-q:
			r.process(ctx, t.Log, t.Description, t.ID, t.Action)
			metrics.PerIDQueue.WithLabelValues(hostID).Dec()
		case <-time.After(r.orchestrator.workerIdleTimeout):
			// TODO: check queue length returning maybe?
			return
		}
	}
}
