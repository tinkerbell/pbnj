package taskrunner

import (
	"context"

	"github.com/adrianbrad/queue"
	"github.com/go-logr/logr"
)

type IngestQueue struct {
	q *queue.Blocking[*Task]
}

type Task struct {
	ID           string                            `json:"id"`
	Host         string                            `json:"host"`
	Description  string                            `json:"description"`
	Action       func(chan string) (string, error) `json:"-"`
	Log          logr.Logger                       `json:"-"`
}

func NewFIFOChannelQueue() *IngestQueue {
	return &IngestQueue{
		q: queue.NewBlocking([]*Task{}),
	}
}

// Enqueue inserts the item into the queue.
func (i *IngestQueue) Enqueue2(item Task) {
	i.q.OfferWait(&item)
}

// Dequeue removes the oldest element from the queue. FIFO.
func (i *IngestQueue) Dequeue2(ctx context.Context, tChan chan Task) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			item, err := i.q.Get()
			if err != nil {
				continue
			}
			tChan <- *item
		}
	}
}

func NewIngestQueue() *IngestQueue {
	return &IngestQueue{
		q: queue.NewBlocking([]*Task{}),
	}
}

// Enqueue inserts the item into the queue.
func (i *IngestQueue) Enqueue(item Task) {
	i.q.OfferWait(&item)
}

// Dequeue removes the oldest element from the queue. FIFO.
func (i *IngestQueue) Dequeue() (Task, error) {
	item, err := i.q.Get()
	if err != nil {
		return Task{}, err
	}

	return *item, nil
}

func (i *IngestQueue) Size() int {
	return i.q.Size()
}

func newHostQueue() *hostQueue {
	return &hostQueue{
		q:  queue.NewBlocking[host]([]host{}),
		ch: make(chan host),
	}
}

type host string

func (h host) String() string {
	return string(h)
}

type hostQueue struct {
	q  *queue.Blocking[host]
	ch chan host
}

// Enqueue inserts the item into the queue.
func (i *hostQueue) Enqueue(item host) {
	i.q.OfferWait(item)
}

// Dequeue removes the oldest element from the queue. FIFO.
func (i *hostQueue) Dequeue() (host, error) {
	item, err := i.q.Get()
	if err != nil {
		return "", err
	}

	return item, nil
}

// Dequeue removes the oldest element from the queue. FIFO.
func (i *hostQueue) Dequeue2(ctx context.Context) <-chan host {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				item, err := i.q.Get()
				if err != nil {
					continue
				}
				i.ch <- item
			}
		}
	}()
	return i.ch
}

func (i *hostQueue) Size() int {
	return i.q.Size()
}
