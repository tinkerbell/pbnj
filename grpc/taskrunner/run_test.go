package taskrunner

import (
	"context"
	"testing"
	"time"

	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	"github.com/tinkerbell/pbnj/grpc/persistence"
)

func TestStart(t *testing.T) {
	t.Skip("skipping test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := freecache.NewStore(freecache.DefaultOptions)
	s := gokv.Store(f)
	repo := &persistence.GoKV{Store: s, Ctx: ctx}

	r := NewRunner(repo, 100, time.Second)
	r.Start(ctx)

	r.orchestrator.manager.WaitAllDone()
	cancel()
	t.Fatal()
}

/*
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
*/
