package runner

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/manifoldco/promptui"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	v1Client "github.com/tinkerbell/pbnj/client"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

type expected struct {
	Action   interface{}
	Want     *v1.StatusResponse
	WaitTime time.Duration
}

type testResource struct {
	Host     string
	Username string
	Password string
	Vendor   string
	Tests    map[string]expected
}

type dataObject map[string]testResource

// RunTests actions against BMCs.
func RunTests(t logr.Logger, cfgData ConfigFile) {
	resources := createTestData(cfgData.Data)
	for _, rs := range resources {
		rs := rs
		tests := rs.Tests
		testsKeys := sortedResources(tests)

		prompt := promptui.Select{
			Label:     fmt.Sprintf("quit, run or skip tests against: %v?", rs.Host),
			Items:     []string{"Run", "Skip", "Quit"},
			Templates: customPromptTemplate(rs.Host),
		}

		_, resp, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		if resp == "Skip" {
			continue
		} else if resp == "Quit" {
			return
		}
		fmt.Println("========", "Start of tests against:", rs.Host, "========")
		for _, key := range testsKeys {
			key := key

			stepName := tests[key]
			t.V(1).Info("debugging", "nextStep", fmt.Sprintf("%v", stepName))
			prompt := promptui.Select{
				Label:     fmt.Sprintf("quit, run or skip step: %v?", key),
				Items:     []string{"Run", "Skip", "Quit"},
				Templates: customPromptTemplate(stepName.Want.Description),
			}

			_, resp, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			if resp == "Skip" {
				continue
			} else if resp == "Quit" {
				return
			}

			fmt.Printf("============ running test: %v ============\n", key)
			whatAmI := func(i interface{}) (*v1.StatusResponse, error) {
				var got *v1.StatusResponse
				var err error
				switch act := i.(type) {
				case v1.PowerAction:
					got, err = runMachinePowerClient(rs, act, cfgData.Server)
					if err != nil {
						return nil, err
					}
				case v1.BootDevice:
					got, err = runMachineBootDevClient(rs, act, cfgData.Server)
					if err != nil {
						return nil, err
					}
				case *v1.CreateUserRequest:
					got, err = runBMCCreateUserClient(rs, act, cfgData.Server)
					if err != nil {
						return nil, err
					}
				case *v1.DeleteUserRequest:
					got, err = runBMCDeleteUserClient(rs, act, cfgData.Server)
					if err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("case not found for %T", act)
				}
				return got, nil
			}
			successful := color.GreenString("PASS")
			var diff string
			got, err := whatAmI(stepName.Action)
			if err != nil {
				log.Println(err)
				successful = color.RedString("FAIL")
				goto COMPLETE
			}
			if got == nil {
				log.Printf("got nil: %v, expected: %v", got, stepName.Want)
				successful = color.RedString("FAIL")
				goto COMPLETE
			}
			got.Result = strings.ToLower(got.Result)
			diff = cmp.Diff(stepName.Want, got, cmpopts.IgnoreMapEntries(func(key, i interface{}) bool { return key == "id" }), protocmp.Transform())
			if diff != "" {
				log.Println(diff)
				successful = color.RedString("FAIL")
			}
		COMPLETE:
			fmt.Printf("============ %v: completed test: %v ============\n", successful, stepName.Want.Description)
		}
		fmt.Println()
		fmt.Println("========", "End of tests against:", rs.Host, "========")
	}
}

func runMachinePowerClient(in testResource, action v1.PowerAction, s Server) (*v1.StatusResponse, error) {
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
		PowerAction: action,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func runMachineBootDevClient(in testResource, action v1.BootDevice, s Server) (*v1.StatusResponse, error) {
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

	inRequest := &v1.DeviceRequest{
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
		BootDevice: action,
		Persistent: false,
		EfiBoot:    false,
	}
	resp, err := v1Client.MachineBootDev(ctx, client, taskClient, inRequest)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func runBMCCreateUserClient(in testResource, action *v1.CreateUserRequest, s Server) (*v1.StatusResponse, error) {
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
	client := v1.NewBMCClient(conn)
	taskClient := v1.NewTaskClient(conn)

	inRequest := &v1.CreateUserRequest{
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
		UserCreds: &v1.UserCreds{
			Username: action.UserCreds.Username,
			Password: action.UserCreds.Password,
			UserRole: action.UserCreds.UserRole,
		},
	}
	resp, err := v1Client.BMCCreateUser(ctx, client, taskClient, inRequest)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func runBMCDeleteUserClient(in testResource, action *v1.DeleteUserRequest, s Server) (*v1.StatusResponse, error) {
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
	client := v1.NewBMCClient(conn)
	taskClient := v1.NewTaskClient(conn)

	inRequest := &v1.DeleteUserRequest{
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
		Username: action.Username,
	}
	resp, err := v1Client.BMCDeleteUser(ctx, client, taskClient, inRequest)
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

func sortedResources(m map[string]expected) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func customPromptTemplate(stepName string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   fmt.Sprintf("%v {{ . }}", promptui.IconSelect),
		Inactive: "  {{ . }}",
		Selected: fmt.Sprintf("%v {{ . | faint }} - %v", promptui.IconGood, promptui.Styler(promptui.FGFaint)(stepName)),
	}
}
