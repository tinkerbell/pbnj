package zaplog

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"
)

func TestRegisterLogger(t *testing.T) {
	t.Skip("will write")
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = writer
	os.Stderr = writer
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
