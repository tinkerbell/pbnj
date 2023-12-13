package rpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bmc-toolbox/bmclib/v2/bmc"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/onsi/gomega"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/freecache"
	"github.com/stretchr/testify/mock"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/grpc/persistence"
	"github.com/tinkerbell/pbnj/grpc/taskrunner"
)

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

func TestDiagnosticService_SystemEventLogRaw(t *testing.T) {
	testCases := []struct {
		name               string
		ctx                context.Context
		in                 *v1.SystemEventLogRawRequest
		expectedLog        *v1.SystemEventLogRawResponse
		expectedErr        error
		expectedLogMessage string
	}{
		{
			name: "successful request",
			ctx:  context.Background(),
			in: &v1.SystemEventLogRawRequest{
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
			expectedLog: &v1.SystemEventLogRawResponse{
				Log: "test log",
			},
			expectedErr:        nil,
			expectedLogMessage: "start Get System Event Log request",
		},
		{
			name: "error creating raw system event log action",
			ctx:  context.Background(),
			in: &v1.SystemEventLogRawRequest{
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
			expectedLog:        nil,
			expectedErr:        errors.New("error creating raw system event log action"),
			expectedLogMessage: "start Get System Event Log request",
		},
		{
			name: "error getting raw system event log",
			ctx:  context.Background(),
			in: &v1.SystemEventLogRawRequest{
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
			expectedLog:        nil,
			expectedErr:        errors.New("error getting raw system event log"),
			expectedLogMessage: "start Get System Event Log request",
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			diagnosticService := NewMockDiagnosticService()
			diagnosticService.On("SystemEventLogRaw", testCase.ctx, testCase.in).Return(testCase.expectedLog, testCase.expectedErr)
		})
	}
}
func TestDiagnosticService_SystemEventLog(t *testing.T) {
	testCases := []struct {
		name           string
		ctx            context.Context
		in             *v1.SystemEventLogRequest
		expectedEvents []*v1.SystemEventLogEntry
		expectedErr    error
	}{
		{
			name: "successful request",
			ctx:  context.Background(),
			in: &v1.SystemEventLogRequest{
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
			ctx:  context.Background(),
			in: &v1.SystemEventLogRequest{
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
			ctx:  context.Background(),
			in: &v1.SystemEventLogRequest{
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
			diagnosticService := NewMockDiagnosticService()
			diagnosticService.On("SystemEventLog", testCase.ctx, testCase.in).Return(&v1.SystemEventLogResponse{
				Events: testCase.expectedEvents,
			}, testCase.expectedErr)

			response, err := diagnosticService.SystemEventLog(testCase.ctx, testCase.in)

			if diff := cmp.Diff(testCase.expectedEvents, response.Events, cmpopts.IgnoreUnexported(v1.SystemEventLogEntry{})); diff != "" {
				t.Errorf("Mismatch in events (-expected, +actual):\n%s", diff)
			}

			if diff := cmp.Diff(testCase.expectedErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Mismatch in error (-expected, +actual):\n%s", diff)
			}
		})
	}
}

type MockDiagnosticService struct {
	mock.Mock
}

func NewMockDiagnosticService() *MockDiagnosticService {
	return &MockDiagnosticService{}
}

func (m *MockDiagnosticService) Screenshot(ctx context.Context, in *v1.ScreenshotRequest) (*v1.ScreenshotResponse, error) {
	args := m.Called(ctx, in)
	resp, ok := args.Get(0).(*v1.ScreenshotResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected type for ScreenshotResponse")
	}
	return resp, args.Error(1)
}

func (m *MockDiagnosticService) SystemEventLog(ctx context.Context, in *v1.SystemEventLogRequest) (*v1.SystemEventLogResponse, error) {
	args := m.Called(ctx, in)
	resp, ok := args.Get(0).(*v1.SystemEventLogResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected type for SystemEventLogResponse")
	}
	return resp, args.Error(1)
}

func (m *MockDiagnosticService) SystemEventLogRaw(ctx context.Context, in *v1.SystemEventLogRawRequest) (*v1.SystemEventLogRawResponse, error) {
	args := m.Called(ctx, in)
	resp, ok := args.Get(0).(*v1.SystemEventLogRawResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected type for SystemEventLogRawResponse")
	}
	return resp, args.Error(1)
}
