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

	ingestManager *concurrencyManager

	fifoQueue      *hostQueue
	fifoChan       chan host
	ingestionQueue *IngestQueue
	// perIDQueue is a map of hostID to a channel of tasks.
	perIDQueue sync.Map

	//testing new stuff
	ingestChan chan Task
}

func (r *Runner) Print() {
	one := r.orchestrator.ingestionQueue.Size()
	two := r.orchestrator.fifoQueue.Size()
	var three int
	r.orchestrator.perIDQueue.Range(func(key, value interface{}) bool {
		three++
		return true
	})
	fmt.Printf("ingestion queue size: %d\n", one)
	fmt.Printf("fcfs queue size: %d\n", two)
	fmt.Printf("perID queue size: %d\n", three)
}

// ingest take a task off the ingestion queue and puts it on the perID queue
// and adds the host ID to the fcfs queue.
func (r *Runner) ingest(ctx context.Context) {
	//func (o *orchestrator) ingest(ctx context.Context) {
	// dequeue from ingestion queue
	// enqueue to perID queue
	// enqueue to fcfs queue
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-r.orchestrator.ingestChan:

			// 2. enqueue to perID queue
			ch := make(chan Task, 5000)
			q, exists := r.orchestrator.perIDQueue.LoadOrStore(t.Host, ch)
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
		time.Sleep(time.Second * 2)
		// r.orchestrator.perIDQueue.Range(func(key, value interface{}) bool {
		r.orchestrator.workers.Range(func(key, value interface{}) bool {
			// if worker id exists in o.workers, then move on because the worker is already running.
			if value.(bool) {
				return true
			}

			// wait for a worker to become available
			r.orchestrator.manager.Wait()

			r.orchestrator.workers.Store(key.(string), true)
			v, _ := r.orchestrator.perIDQueue.Load(key.(string))
			go r.worker(ctx, key.(string), v.(chan Task))
			return true
		})
	}
}

func (r *Runner) worker(ctx context.Context, hostID string, q chan Task) {
	defer r.orchestrator.manager.Done()
	defer func() {
		r.orchestrator.workers.Range(func(key, value interface{}) bool {
			if key.(string) == hostID { //nolint:forcetypeassert // good
				r.orchestrator.workers.Delete(key.(string))
				return true //nolint:revive // this is needed to satisfy the func parameter
			}
			return true //nolint:revive // this is needed to satisfy the func parameter
		})

	}()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-q:
			r.process(ctx, t.Log, t.Description, t.ID, t.Action)
			metrics.PerIDQueue.WithLabelValues(hostID).Dec()
		case <-time.After(r.orchestrator.workerIdleTimeout):
			return
		}
	}
}
