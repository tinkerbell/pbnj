package taskrunner

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	o := Start(ctx)

	o.manager.WaitAllDone()

	
	time.Sleep(3 * time.Second)
	t.Log("=======", runtime.NumGoroutine())
}
