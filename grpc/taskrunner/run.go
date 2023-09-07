package taskrunner

import (
	"context"
	"sync"
	"time"

	"github.com/tinkerbell/pbnj/pkg/metrics"
)

type orchestrator struct {
	workers           sync.Map
	manager           *concurrencyManager
	workerIdleTimeout time.Duration

	ingestManager *concurrencyManager

	fifoQueue      *hostQueue
	ingestionQueue *IngestQueue
	// perIDQueue is a map of hostID to a channel of tasks.
	perIDQueue sync.Map
}

// ingest take a task off the ingestion queue and puts it on the perID queue
// and adds the host ID to the fcfs queue.
func (o *orchestrator) ingest(ctx context.Context) {
	// dequeue from ingestion queue
	// enqueue to perID queue
	// enqueue to fcfs queue
	for {
		select {
		case <-ctx.Done():
			return
		default:
			o.ingestManager.Wait()
			go func() {
				defer o.ingestManager.Done()
				// 1. dequeue from ingestion queue
				// 2. enqueue to perID queue
				// 3. enqueue to fcfs queue
				// ---
				// 1. dequeue from ingestion queue
				t, err := o.ingestionQueue.Dequeue()
				if err != nil {
					return
				}
				metrics.IngestionQueue.Dec()

				// 2. enqueue to perID queue
				que := NewIngestQueue()
				q, _ := o.perIDQueue.LoadOrStore(t.Host, que)
				v, ok := q.(*IngestQueue)
				if !ok {
					return
				}
				v.Enqueue(t)
				metrics.PerIDQueue.WithLabelValues(t.Host).Inc()

				// 3. enqueue to fcfs queue
				o.fifoQueue.Enqueue(host(t.Host))
			}()
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
		select {
		case <-ctx.Done():
			return
		default:
			// 1. dequeue from fcfs queue
			// 2. if worker id exists in o.workers, then move on
			// 2a. if len(o.workers) < o.maxWorkers then create worker and move on, else recursion.
			h, err := r.orchestrator.fifoQueue.Dequeue()
			if err != nil {
				continue
			}

			// check queue length for perID queue, if 0, then continue
			if elem, ok := r.orchestrator.perIDQueue.Load(h.String()); ok {
				v, ok := elem.(*IngestQueue)
				if !ok {
					continue
				}
				if v.Size() == 0 {
					continue
				}
			}

			// if worker id exists in o.workers, then move on because the worker is already running.
			_, ok := r.orchestrator.workers.Load(h)
			if ok {
				continue
			}

			// wait for a worker to become available
			r.orchestrator.manager.Wait()

			r.orchestrator.workers.Store(h, true)
			go r.worker(ctx, h.String())
		}
	}
}

func (r *Runner) worker(ctx context.Context, hostID string) {
	defer r.orchestrator.manager.Done()
	defer func() {
		r.orchestrator.workers.Range(func(key, value interface{}) bool {
			if key.(host) == host(hostID) { //nolint:forcetypeassert // good
				r.orchestrator.workers.Delete(key)
				return true //nolint:revive // this is needed to satisfy the func parameter
			}
			return true //nolint:revive // this is needed to satisfy the func parameter
		})
	}()

	elem, ok := r.orchestrator.perIDQueue.Load(hostID)
	if !ok {
		return
	}

	tChan := make(chan Task)
	cctx, stop := context.WithCancel(ctx)
	defer stop()

	go func(ch chan Task) {
		defer close(ch)
		for {
			select {
			case <-cctx.Done():
				return
			default:
				task, err := elem.(*IngestQueue).Dequeue()
				if err != nil {
					continue
				}
				tChan <- task
				metrics.PerIDQueue.WithLabelValues(hostID).Dec()
			}
		}
	}(tChan)

	for {
		select {
		case <-ctx.Done():
			stop()
			return
		case t := <-tChan:
			r.process(ctx, t.Log, t.Description, t.ID, t.Action)
		case <-time.After(r.orchestrator.workerIdleTimeout):
			stop()
			return
		}
	}
}
