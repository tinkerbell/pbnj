// +build functional

package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	v1Client "github.com/tinkerbell/pbnj/client"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	lookup = map[string]map[string]expected{
		"happyTests":           happyTests,
		"notIdentifiableTests": notIdentifiableTests,
		"happyTestsOneOff":     happyTestsOneOff,
	}
	happyTests = map[string]expected{
		"power status": {
			Action: v1.PowerRequest_STATUS,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "power action",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "on",
				Complete:    true,
				Messages:    []string{"trying to connect to bmc", "connected to bmc", "getting power status"},
			},
		},
		"power on":      {Action: v1.PowerRequest_ON, Want: notImplementedWant("ON")},
		"power off":     {Action: v1.PowerRequest_OFF, Want: notImplementedWant("OFF")},
		"power hardoff": {Action: v1.PowerRequest_HARDOFF, Want: notImplementedWant("HARD OFF")},
		"power cycle":   {Action: v1.PowerRequest_CYCLE, Want: notImplementedWant("CYCLE")},
		"power reset":   {Action: v1.PowerRequest_RESET, Want: notImplementedWant("RESET")},
	}
	notIdentifiableTests = map[string]expected{
		"power status":  {Action: v1.PowerRequest_STATUS, Want: notIdentifiableWant},
		"power on":      {Action: v1.PowerRequest_ON, Want: notIdentifiableWant},
		"power off":     {Action: v1.PowerRequest_OFF, Want: notIdentifiableWant},
		"power hardoff": {Action: v1.PowerRequest_HARDOFF, Want: notIdentifiableWant},
		"power cycle":   {Action: v1.PowerRequest_CYCLE, Want: notIdentifiableWant},
		"power reset":   {Action: v1.PowerRequest_RESET, Want: notIdentifiableWant},
	}
	notIdentifiableWant = &v1.StatusResponse{
		Id:          "12345",
		Description: "power action",
		Error: &v1.Error{
			Code:    2,
			Message: "unable to identify the vendor",
			Details: nil,
		},
		State:    "complete",
		Result:   "action failed",
		Complete: true,
		Messages: []string{"trying to connect to bmc"},
	}
	happyTestsOneOff = updateSingleTest("power status", happyTests, expected{
		Action: v1.PowerRequest_STATUS,
		Want: &v1.StatusResponse{
			Id:          "12345",
			Description: "power action",
			Error: &v1.Error{
				Code:    2,
				Message: "XML syntax error on line 10: element <META> closed by </head>",
				Details: nil,
			},
			State:    "complete",
			Result:   "action failed",
			Complete: true,
			Messages: []string{"trying to connect to bmc", "connected to bmc", "getting power status", "error getting power state"},
		},
	})
)

type expected struct {
	Action v1.PowerRequest_Action
	Want   *v1.StatusResponse
}

type testResource struct {
	Host     string
	Username string
	Password string
	Vendor   string
	Tests    map[string]expected
}

type dataObject map[string]testResource

// TestPower actions against BMCs
func TestPower(t *testing.T) {
	resources := createTestData(cfgData.Data)
	for rname, rs := range resources {
		rs := rs
		rname := rname
		t.Run(rname, func(t *testing.T) {
			t.Parallel()
			for name, tc := range rs.Tests {
				tc := tc
				name := name
				t.Run(name, func(t *testing.T) {
					// do the work
					got, err := runMachineClient(rs, tc.Action, cfgData.Server)
					if err != nil {
						t.Fatal(err)
					}

					got.Id = "12345"
					diff := cmp.Diff(tc.Want, got, protocmp.Transform())
					if diff != "" {
						t.Fatalf(diff)
					}
				})
			}
		})
	}
}

func runMachineClient(in testResource, action v1.PowerRequest_Action, s Server) (*v1.StatusResponse, error) {
	var opts []grpc.DialOption
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(s.URL+":"+s.Port, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := v1.NewMachineClient(conn)
	taskClient := v1.NewTaskClient(conn)

	resp, err := v1Client.MachinePower(ctx, client, taskClient, &v1.PowerRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: &v1.DirectAuthn{
					Host: &v1.Host{
						Host: in.Host,
					},
					Username: in.Username,
					Password: in.Password,
				},
			},
		},
		Vendor: &v1.Vendor{
			Name: in.Vendor,
		},
		Action: action,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func createTestData(config Data) dataObject {
	f := dataObject{}
	for _, elem := range config.Resources {
		tmp := testResource{}
		tmp.Host = elem.IP
		tmp.Password = elem.Password
		tmp.Username = elem.Username
		tmp.Vendor = elem.Vendor

		tests := map[string]expected{}
		for _, elem := range elem.UseCases.Power {
			t, ok := lookup[elem]
			if !ok {
				fmt.Printf("FYI, useCase: '%v' was not found. please verify it exists in the code base\n", elem)
			}
			for k, v := range t {
				tests[elem+"_"+k] = v
			}
		}
		tmp.Tests = tests

		f[elem.IP] = tmp
	}
	return f
}

func notImplementedWant(fn string) *v1.StatusResponse {
	return &v1.StatusResponse{
		Id:          "12345",
		Description: "power action",
		Error: &v1.Error{
			Code:    12,
			Message: fmt.Sprintf("power %v not implemented", fn),
			Details: nil,
		},
		State:    "complete",
		Result:   "action failed",
		Complete: true,
		Messages: []string{"trying to connect to bmc", "connected to bmc"},
	}
}

func updateSingleTest(key string, existing map[string]expected, val expected) map[string]expected {
	newExisting := make(map[string]expected)
	for k, v := range existing {
		newExisting[k] = v
	}
	newExisting[key] = val
	return newExisting
}
