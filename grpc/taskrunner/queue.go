package taskrunner

import (
	"github.com/adrianbrad/queue"
	"github.com/go-logr/logr"
)

type IngestQueue struct {
	q *queue.Blocking[*Task]
}

type Task struct {
	ID          string                            `json:"id"`
	Host        string                            `json:"host"`
	Description string                            `json:"description"`
	Action      func(chan string) (string, error) `json:"-"`
	Log         logr.Logger                       `json:"-"`
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
		q: queue.NewBlocking[host]([]host{}, queue.WithCapacity(10000)),
	}
}

type host string

func (h host) String() string {
	return string(h)
}

type hostQueue struct {
	q *queue.Blocking[host]
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

func (i *hostQueue) Size() int {
	return i.q.Size()
}
