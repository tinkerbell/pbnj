package oob

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/packethost/pkg/log/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

func TestParseAuth(t *testing.T) {

	tests := map[string]struct {
		input *v1.Authn
		want  repository.Error
	}{
		"with auth":       {input: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{Host: &v1.Host{Host: "localhost"}, Username: "admin", Password: "admin"}}}, want: repository.Error{}},
		"nil Direct Auth": {input: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: nil}}, want: repository.Error{Code: v1.Code_value["UNAUTHENTICATED"], Message: "no auth found", Details: nil}},
		"nil auth":        {input: nil, want: repository.Error{Code: v1.Code_value["UNAUTHENTICATED"], Message: "no auth found", Details: nil}},
	}
	ctx := context.Background()
	l, zapLogger, _ := logr.NewPacketLogr()
	ctx = ctxzap.ToContext(ctx, zapLogger)
	sm := make(chan string)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			a := Accessory{
				Log:            l,
				Ctx:            ctx,
				StatusMessages: sm,
			}

			host, username, passwd, errMsg := a.ParseAuth(tc.input)
			diff := cmp.Diff(tc.want, errMsg)
			if diff != "" {
				t.Log(fmt.Sprintf("%+v", errMsg))
				t.Fatalf(diff)
			}

			expectedHost := tc.input.GetDirectAuthn().GetHost().GetHost()
			expectedUsername := tc.input.GetDirectAuthn().GetUsername()
			expectedPassword := tc.input.GetDirectAuthn().GetPassword()

			if host != expectedHost {
				t.Fatalf("expectedHost: %v, got: %v", expectedHost, host)
			}
			if username != expectedUsername {
				t.Fatalf("expectedUsername: %v, got: %v", expectedUsername, username)
			}
			if passwd != expectedPassword {
				t.Fatalf("expectedPassword: %v, got: %v", expectedPassword, passwd)
			}
		})
	}
}

func TestSendStatusMessage(t *testing.T) {

	tests := map[string]struct {
		runChanReceiver bool
		want            []string
	}{
		"with chan receiver":    {runChanReceiver: true, want: []string{"test message 1", "test message 2"}},
		"without chan receiver": {runChanReceiver: false, want: nil},
	}

	ctx := context.Background()
	l, zapLogger, _ := logr.NewPacketLogr()
	ctx = ctxzap.ToContext(ctx, zapLogger)
	sm := make(chan string)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var msgs []string
			done := make(chan bool, 1)
			a := Accessory{
				Log:            l,
				Ctx:            ctx,
				StatusMessages: sm,
			}

			if tc.runChanReceiver {
				go func() {
					for {
						msgs = append(msgs, <-a.StatusMessages)
						select {
						case <-done:
							return
						default:
							continue
						}
					}
				}()
				a.SendStatusMessage(tc.want[0])
				a.SendStatusMessage(tc.want[1])
				time.Sleep(10 * time.Millisecond)
				done <- true
			}

			diff := cmp.Diff(msgs, tc.want)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}

}
