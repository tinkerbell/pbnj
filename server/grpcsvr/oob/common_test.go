package oob

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/packethost/pkg/log/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

func TestParseAuth(t *testing.T) {
	tests := map[string]struct {
		input *v1.Authn
		want  *repository.Error
	}{
		"with auth":       {input: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{Host: &v1.Host{Host: "localhost"}, Username: "admin", Password: "admin"}}}, want: nil},
		"nil Direct Auth": {input: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: nil}}, want: &repository.Error{Code: v1.Code_value["UNAUTHENTICATED"], Message: "no auth found", Details: nil}},
		"nil auth":        {input: nil, want: &repository.Error{Code: v1.Code_value["UNAUTHENTICATED"], Message: "no auth found", Details: nil}},
	}
	packetLogr, _, _ := logr.NewPacketLogr()
	sm := make(chan string)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			a := Accessory{
				Log:            packetLogr.Logger,
				StatusMessages: sm,
			}

			host, username, passwd, errMsg := a.ParseAuth(tc.input)
			if errMsg != nil {
				diff := cmp.Diff(tc.want.Error(), errMsg.Error())
				if diff != "" {
					t.Log(fmt.Sprintf("%+v", errMsg))
					t.Fatalf(diff)
				}
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

	packetLogr, _, _ := logr.NewPacketLogr()
	sm := make(chan string)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var msgs []string
			done := make(chan bool, 1)
			a := Accessory{
				Log:            packetLogr.Logger,
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

type notAConnectionInterface struct{}

type testConnect struct {
	makeFail bool
}

func (t *testConnect) Connect(ctx context.Context) error {
	if t.makeFail {
		return errors.New("failed to connect")
	}
	return nil
}

func (t *testConnect) Close(ctx context.Context) {}

func TestEstablishConnections(t *testing.T) {
	testCases := []struct {
		name                  string
		makeFail              bool
		successfulConnections []string
		implementation        *testConnect
		notAnImplementation   *notAConnectionInterface
		err                   error
	}{
		{name: "success", successfulConnections: []string{"test"}, implementation: &testConnect{}},
		{name: "fail to connect", makeFail: true, implementation: &testConnect{makeFail: true}, err: errors.New("code: 2 message: could not connect details: [failed to connect]")},
		{name: "unknown connection type", makeFail: true, notAnImplementation: &notAConnectionInterface{}, err: errors.New("code: 2 message: could not connect details: [unknown connection type]")},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var testImplementation map[string]interface{}
			if tc.implementation != nil {
				testImplementation = map[string]interface{}{"test": tc.implementation}
			} else {
				testImplementation = map[string]interface{}{"test": tc.notAnImplementation}
			}

			result, err := EstablishConnections(context.Background(), testImplementation)
			if err != nil {
				t.Log(err)
				diff := cmp.Diff(tc.err.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}

			} else {
				diff := cmp.Diff(tc.successfulConnections, result)
				if diff != "" {
					t.Fatal(diff)
				}
				Close(context.Background(), tc.implementation)
			}

		})
	}
}
