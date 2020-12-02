package oob

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-multierror"
)

type OOBTester struct {
	MakeFail bool
}

func (o *OOBTester) PowerSet(ctx context.Context, action string) (result string, err error) {
	if o.MakeFail {
		return result, errors.New("power failed")
	}
	return "power action complete: " + action, nil
}

func (o *OOBTester) BootDeviceSet(ctx context.Context, device string) (result string, err error) {
	if o.MakeFail {
		return result, errors.New("boot device failed")
	}
	return "boot device set: " + device, nil
}

func (o *OOBTester) BMCReset(ctx context.Context, rType string) (err error) {
	if o.MakeFail {
		return errors.New("failed: BMC reset")
	}
	return nil
}

func (o *OOBTester) CreateUser(ctx context.Context) (err error) {
	if o.MakeFail {
		return errors.New("create user failed")
	}
	return nil
}

func (o *OOBTester) UpdateUser(ctx context.Context) (err error) {
	if o.MakeFail {
		return errors.New("update user failed")
	}
	return nil
}

func (o *OOBTester) DeleteUser(ctx context.Context) (err error) {
	if o.MakeFail {
		return errors.New("delete user failed")
	}
	return nil
}

func TestMachinePower(t *testing.T) {
	testCases := []struct {
		name     string
		action   string
		makeFail bool
		err      error
	}{
		{name: "success", action: "status", err: nil},
		{name: "Power method fails", action: "status", makeFail: true, err: &multierror.Error{Errors: []error{errors.New("power failed"), errors.New("power state failed")}}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testImplementation := OOBTester{MakeFail: tc.makeFail}
			expectedResult := "power action complete: " + tc.action
			result, err := SetPower(context.Background(), tc.action, []PowerSetter{&testImplementation})
			if err != nil {
				diff := cmp.Diff(tc.err.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}

			} else {
				diff := cmp.Diff(expectedResult, result)
				if diff != "" {
					t.Fatal(diff)
				}
			}

		})
	}
}

func TestMachineBootDevice(t *testing.T) {
	testCases := []struct {
		name     string
		device   string
		makeFail bool
		err      error
	}{
		{name: "success", device: "pxe", err: nil},
		{name: "Power method fails", device: "pxe", makeFail: true, err: &multierror.Error{Errors: []error{errors.New("boot device failed"), errors.New("set boot device failed")}}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testImplementation := OOBTester{MakeFail: tc.makeFail}
			expectedResult := "boot device set: " + tc.device
			result, err := SetBootDevice(context.Background(), tc.device, []BootDeviceSetter{&testImplementation})
			if err != nil {
				diff := cmp.Diff(tc.err.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}

			} else {
				diff := cmp.Diff(expectedResult, result)
				if diff != "" {
					t.Fatal(diff)
				}
			}

		})
	}
}

func TestBMCReset(t *testing.T) {
	testCases := []struct {
		name      string
		resetType string
		makeFail  bool
		err       error
	}{
		{name: "success", resetType: "cold", err: nil},
		{name: "BMC reset fails", resetType: "cold", makeFail: true, err: &multierror.Error{Errors: []error{errors.New("failed: BMC reset"), errors.New("BMC reset failed")}}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testImplementation := OOBTester{MakeFail: tc.makeFail}
			err := ResetBMC(context.Background(), tc.resetType, []BMCResetter{&testImplementation})
			if err != nil {
				diff := cmp.Diff(tc.err.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name     string
		makeFail bool
		err      error
	}{
		{name: "success", err: nil},
		{name: "Create User fails", makeFail: true, err: &multierror.Error{Errors: []error{errors.New("create user failed"), errors.New("create user failed")}}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testImplementation := OOBTester{MakeFail: tc.makeFail}
			err := CreateUser(context.Background(), []BMC{&testImplementation})
			if err != nil {
				diff := cmp.Diff(tc.err.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name     string
		makeFail bool
		err      error
	}{
		{name: "success", err: nil},
		{name: "Update User fails", makeFail: true, err: &multierror.Error{Errors: []error{errors.New("update user failed"), errors.New("update user failed")}}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testImplementation := OOBTester{MakeFail: tc.makeFail}
			err := UpdateUser(context.Background(), []BMC{&testImplementation})
			if err != nil {
				diff := cmp.Diff(tc.err.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name     string
		makeFail bool
		err      error
	}{
		{name: "success", err: nil},
		{name: "Delete User fails", makeFail: true, err: &multierror.Error{Errors: []error{errors.New("delete user failed"), errors.New("delete user failed")}}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testImplementation := OOBTester{MakeFail: tc.makeFail}
			err := DeleteUser(context.Background(), []BMC{&testImplementation})
			if err != nil {
				diff := cmp.Diff(tc.err.Error(), err.Error())
				if diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}