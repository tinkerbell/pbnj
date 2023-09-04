package taskrunner

import (
	"context"
	"testing"
)

func TestStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	o := Start(ctx)

	o.manager.WaitAllDone()
	cancel()
}
