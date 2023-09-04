package taskrunner

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/gosuri/uilive"
)

type orchestrator struct {
	workers sync.Map
	manager *concurrencyManager

	fifoQueue      *hostQueue
	ingestionQueue *IngestQueue
	// perIDQueue is a map of hostID to a channel of tasks.
	perIDQueue sync.Map
}

func Start(ctx context.Context) *orchestrator {
	o := &orchestrator{
		fifoQueue:      NewHostQueue(),
		ingestionQueue: NewIngestQueue(),
		// perIDQueue is a map of hostID to a channel of tasks.
		perIDQueue: sync.Map{},
		manager:    New(17),
	}
	// 1. start the ingestor
	// 2. start the orchestrator
	go o.ingest(ctx)
	go o.orchestrate(ctx)
	go o.observe(ctx)
	go o.load()

	return o
}

func lenSyncMap(m *sync.Map) int {
	var i int
	m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

func (o *orchestrator) observe(ctx context.Context) {
	writer := uilive.New()
	//writer.Out = fileWriter()

	// start listening for updates and render
	writer.Start()
	defer writer.Stop()
	workers := func() string {
		s := []string{}
		for _, h := range hostList {
			v, ok := o.perIDQueue.Load(h)
			if ok {
				s = append(s, fmt.Sprintf("worker (%v) queue len: %v", h, v.(*IngestQueue).Size()))
			}
		}
		return strings.Join(s, "\n") + "\n"
	}
	for {
		select {
		case <-ctx.Done():
			writer.Flush()
			return
		default:
			one := fmt.Sprintf("running workers count: %v\n", o.manager.RunningCount())
			two := fmt.Sprintf("fifoQueue queue length: %v\n", o.fifoQueue.q.Size())
			three := fmt.Sprintf("goroutines: %v\n", runtime.NumGoroutine())
			four := fmt.Sprintf("worker map: %v\n", lenSyncMap(&o.workers))
			all := fmt.Sprintf("%v%v%v%v%v", one, two, three, four, workers())
			fmt.Fprint(writer, all)
			<-time.After(time.Millisecond * 50)
		}
	}
}

func (o *orchestrator) load() {
	for i := 0; i < 1000; i++ {
		i := i
		go func() {
			o.ingestionQueue.Enqueue(Task{
				ID:          strconv.Itoa(i),
				Host:        hostList[rand.Intn(len(hostList))],
				Description: "",
				Action: func(chan string) (string, error) {
					time.Sleep(time.Millisecond * 250)
					return "", nil
				},
				Log: logr.Discard(),
			})
		}()
	}
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
			// 1. dequeue from ingestion queue
			// 2. enqueue to perID queue
			// 3. enqueue to fcfs queue
			// ---
			// 1. dequeue from ingestion queue
			t, err := o.ingestionQueue.Dequeue()
			if err != nil {
				continue
			}

			// 2. enqueue to perID queue
			// hostCh := make(chan Task)
			que := NewIngestQueue()
			q, _ := o.perIDQueue.LoadOrStore(t.Host, que)
			q.(*IngestQueue).Enqueue(t)

			// 3. enqueue to fcfs queue
			o.fifoQueue.Enqueue(host(t.Host))
		}
	}
}

// 1. dequeue from fcfs queue
// 2. dequeue from perID queue
// 3. if worker id exists, send task to worker. else continue.
// 4. if maxWorkers is reached, wait for available worker. else create worker and send task to worker.
func (o *orchestrator) orchestrate(ctx context.Context) {
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
			h, err := o.fifoQueue.Dequeue()
			if err != nil {
				continue
			}

			// check queue length for perID queue, if 0, then continue
			if elem, ok := o.perIDQueue.Load(h.String()); ok {
				if elem.(*IngestQueue).Size() == 0 {
					continue
				}
			}

			// if worker id exists in o.workers, then move on because the worker is already running.
			_, ok := o.workers.Load(h)
			if ok {
				continue
			}

			// wait for a worker to become available
			o.manager.Wait()

			o.workers.Store(h, true)
			go o.worker(ctx, h.String())
		}
	}
}

func (o *orchestrator) worker(ctx context.Context, hostID string) {
	defer o.manager.Done()
	defer func() {
		o.workers.Range(func(key, value interface{}) bool {
			if key.(host) == host(hostID) {
				o.workers.Delete(key)
				return true
			}
			return true
		})
	}()

	elem, ok := o.perIDQueue.Load(hostID)
	if !ok {
		return
	}

	tChan := make(chan Task)
	cctx, stop := context.WithCancel(ctx)
	defer stop()
	defer close(tChan)

	go func() {
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
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			stop()
			return
		case t := <-tChan:
			msgCh := make(chan string)
			t.Action(msgCh)
		case <-time.After(2 * time.Second):
			stop()
			return
		}
	}

}

var hostList = []string{
	"127.0.0.1",
	"127.0.0.2",
	"127.0.0.3",
	"127.0.0.4",
	"127.0.0.5",
	"127.0.0.6",
	"127.0.0.7",
	"127.0.0.8",
	"127.0.0.9",
	"127.0.0.10",
	"127.0.0.11",
	"127.0.0.12",
	"127.0.0.13",
	"127.0.0.14",
	"127.0.0.15",
	"127.0.0.16",
	"127.0.0.17",
	"127.0.0.18",
	"127.0.0.19",
	"127.0.0.20",
	"127.0.0.21",
	"127.0.0.22",
	"127.0.0.23",
	"127.0.0.24",
	"127.0.0.25",
	"127.0.0.26",
	"127.0.0.27",
	"127.0.0.28",
	"127.0.0.29",
	"127.0.0.30",
	"127.0.0.31",
	"127.0.0.32",
}
