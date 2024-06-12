package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/bmc-toolbox/bmclib/v2/bmc"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/onsi/gomega"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/persistence"
	"github.com/tinkerbell/pbnj/grpc/taskrunner"
)

type MockDiagnosticAction struct {
	SystemEventLogFunc func(ctx context.Context) (entries []*v1.SystemEventLogEntry, raw string, err error)
}

func (m *MockDiagnosticAction) SystemEventLog(ctx context.Context) (entries []*v1.SystemEventLogEntry, raw string, err error) {
	if m.SystemEventLogFunc != nil {
		return m.SystemEventLogFunc(ctx)
	}
	return nil, "", nil
}

func TestClearSystemEventLog(t *testing.T) {
	testCases := []struct {
		name        string
		req         *v1.ClearSystemEventLogRequest
		expectedErr error
	}{
		{
			name: "status good; direct auth",
			req: &v1.ClearSystemEventLogRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Host: &v1.Host{
								Host: "127.0.0.1",
							},
							Username: "ADMIN",
							Password: "ADMIN",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "",
				},
			},
		},
		{
			name:        "validation failure",
			req:         &v1.ClearSystemEventLogRequest{Authn: &v1.Authn{Authn: &v1.Authn_DirectAuthn{DirectAuthn: &v1.DirectAuthn{}}}},
			expectedErr: errors.New("input arguments are invalid: invalid field Authn.DirectAuthn.Username: value '' must not be an empty string"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)

			ctx := context.Background()

			f := freecache.NewStore(freecache.DefaultOptions)
			s := gokv.Store(f)
			repo := &persistence.GoKV{
				Store: s,
				Ctx:   ctx,
			}

			taskRunner := &taskrunner.Runner{
				Repository: repo,
				Ctx:        ctx,
			}

			diagnosticService := DiagnosticService{
				TaskRunner: taskRunner,
			}

			response, err := diagnosticService.ClearSystemEventLog(ctx, testCase.req)

			t.Log("Got : ", response)
			if err != nil {
				diff := cmp.Diff(testCase.expectedErr.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			} else {
				g.Expect(response.TaskId).Should(gomega.HaveLen(20))
			}
		})
	}
}
func TestConvertEntriesToEvents(t *testing.T) {
	testCases := []struct {
		name           string
		entries        bmc.SystemEventLogEntries
		expectedEvents []*v1.SystemEventLogEntry
	}{
		{
			name:           "empty entries",
			entries:        bmc.SystemEventLogEntries{},
			expectedEvents: nil,
		},
		{
			name: "valid entries",
			entries: bmc.SystemEventLogEntries{
				{"1", "2022-01-01", "Event 1", "Message 1"},
				{"2", "2022-01-02", "Event 2", "Message 2"},
			},
			expectedEvents: []*v1.SystemEventLogEntry{
				{
					Id:          "1",
					Timestamp:   "2022-01-01",
					Description: "Event 1",
					Message:     "Message 1",
				},
				{
					Id:          "2",
					Timestamp:   "2022-01-02",
					Description: "Event 2",
					Message:     "Message 2",
				},
			},
		},
		{
			name: "invalid entries",
			entries: bmc.SystemEventLogEntries{
				{"1", "2022-01-01", "Event 1"},
				{"2", "2022-01-02", "Event 2", "Message 2"},
			},
			expectedEvents: []*v1.SystemEventLogEntry{
				{
					Id:          "2",
					Timestamp:   "2022-01-02",
					Description: "Event 2",
					Message:     "Message 2",
				},
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			events := convertEntriesToEvents(testCase.entries)

			if diff := cmp.Diff(testCase.expectedEvents, events, cmpopts.IgnoreUnexported(v1.SystemEventLogEntry{})); diff != "" {
				t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
			}
		})
	}
}
func TestSystemEventLog(t *testing.T) {
	testCases := []struct {
		name           string
		req            *v1.SystemEventLogRequest
		expectedEvents []*v1.SystemEventLogEntry
		expectedErr    error
	}{
		{
			name: "success",
			req: &v1.SystemEventLogRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Username: "testuser",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "testvendor",
				},
			},
			expectedEvents: []*v1.SystemEventLogEntry{
				{
					Id:          "1",
					Timestamp:   "2022-01-01",
					Description: "Event 1",
					Message:     "Message 1",
				},
				{
					Id:          "2",
					Timestamp:   "2022-01-02",
					Description: "Event 2",
					Message:     "Message 2",
				},
			},
			expectedErr: nil,
		},
		{
			name: "error creating system event log action",
			req: &v1.SystemEventLogRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Username: "testuser",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "testvendor",
				},
			},
			expectedEvents: nil,
			expectedErr:    errors.New("error creating system event log action"),
		},
		{
			name: "error getting system event log",
			req: &v1.SystemEventLogRequest{
				Authn: &v1.Authn{
					Authn: &v1.Authn_DirectAuthn{
						DirectAuthn: &v1.DirectAuthn{
							Username: "testuser",
						},
					},
				},
				Vendor: &v1.Vendor{
					Name: "testvendor",
				},
			},
			expectedEvents: nil,
			expectedErr:    errors.New("error getting system event log"),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			diagnosticService := &MockDiagnosticAction{
				SystemEventLogFunc: func(ctx context.Context) (entries []*v1.SystemEventLogEntry, raw string, err error) {
					if testCase.expectedErr != nil {
						return nil, "", testCase.expectedErr
					}
					return []*v1.SystemEventLogEntry{
						{Id: "1", Timestamp: "2022-01-01", Description: "Event 1", Message: "Message 1"},
						{Id: "2", Timestamp: "2022-01-02", Description: "Event 2", Message: "Message 2"},
					}, "", nil
				},
			}

			response, _, err := diagnosticService.SystemEventLog(ctx)

			t.Log("Got : ", response)
			if err != nil {
				if testCase.expectedErr == nil {
					t.Fatalf("Unexpected error: %v", err)
				} else {
					if diff := cmp.Diff(testCase.expectedErr.Error(), err.Error()); diff != "" {
						t.Fatalf("Error mismatch (-expected, +actual):\n%s", diff)
					}
				}
			} else {
				if diff := cmp.Diff(testCase.expectedEvents, response, cmpopts.IgnoreUnexported(v1.SystemEventLogEntry{})); diff != "" {
					t.Fatalf("Events mismatch (-expected, +actual):\n%s", diff)
				}
			}
		})
	}
}
