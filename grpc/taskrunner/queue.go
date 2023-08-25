package taskrunner

import (
	"sync"

	"github.com/pkg/errors"
)

type IngestQueue struct {
	mu    sync.Mutex
	queue []Ingest
}

type Ingest struct {
	ID          string
	Host        string
	Description string
	Action      func(chan string) (string, error)
}

func NewIngestQueue() *IngestQueue {
	return &IngestQueue{
		queue: []Ingest{},
	}
}

// Enqueue inserts the item into the queue.
func (i *IngestQueue) Enqueue(item Ingest) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.queue = append(i.queue, item)
}

// Dequeue removes the oldest element from the queue. FIFO.
func (i *IngestQueue) Dequeue() (Ingest, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if len(i.queue) > 0 {
		item := i.queue[0]
		i.queue = i.queue[1:]
		return item, nil
	}
	return Ingest{}, errors.New("queue is empty")
}
