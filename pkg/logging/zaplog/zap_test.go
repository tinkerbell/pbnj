package zaplog

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap/zapcore"
)

func TestRegisterLogger(t *testing.T) {
	rLogger, zapLogger, err := RegisterLogger(
		WithKeysAndValues(map[string]interface{}{"hello": "world"}),
		WithLogLevel("debug"),
		WithOutputPaths([]string{"stdout"}),
		WithServiceName("testing"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !rLogger.Enabled() {
		t.Fatal("logger not enabled")
	}
	check := zapLogger.Check(zapcore.InfoLevel, "test message")
	if check.Message != "test message" {
		t.Fatal("error with zap logger")
	}
}

func TestGetContextLogger(t *testing.T) {
	expectedLogMsg := "new logger test message"
	testCtx := context.Background()

	capturedOutput := captureOutput(func() {
		rLogger, zapLogger, err := RegisterLogger(
			WithOutputPaths([]string{"stdout"}),
			WithEnableErrLogsToStderr(true),
		)
		if err != nil {
			t.Fatal(err)
		}
		newCtx := ctxzap.ToContext(testCtx, zapLogger)
		newRLogger := rLogger.GetContextLogger(newCtx)
		if !newRLogger.Enabled() {
			t.Fatal("logger not enabled")
		}
		newRLogger.Info(expectedLogMsg)
	})

	if !strings.Contains(capturedOutput, expectedLogMsg) {
		t.Fatalf("expected to contain: %v, got: %v", expectedLogMsg, capturedOutput)
	}
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
