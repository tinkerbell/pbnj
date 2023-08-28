package taskrunner

import (
	"github.com/adrianbrad/queue"
	"github.com/go-logr/logr"
)

type IngestQueue struct {
	q *queue.Blocking[*Ingest]
}

type Ingest struct {
	ID          string                            `json:"id"`
	Host        string                            `json:"host"`
	Description string                            `json:"description"`
	Action      func(chan string) (string, error) `json:"-"`
	Log         logr.Logger                       `json:"-"`
}

func NewIngestQueue() *IngestQueue {
	return &IngestQueue{
		q: queue.NewBlocking([]*Ingest{}, queue.WithCapacity(10000)),
	}
}

// Enqueue inserts the item into the queue.
func (i *IngestQueue) Enqueue(item Ingest) {
	i.q.OfferWait(&item)
}

// Dequeue removes the oldest element from the queue. FIFO.
func (i *IngestQueue) Dequeue() (Ingest, error) {
	item := i.q.GetWait()

	return *item, nil
}
