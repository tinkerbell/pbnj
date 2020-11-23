package oob

import (
	"context"
	"errors"

	"github.com/hashicorp/go-multierror"
)

// PowerSetter management methods
type PowerSetter interface {
	// Power get status and sets power states like on/off/etc
	PowerSet(ctx context.Context, action string) (result string, err error)
}

// BootDeviceSetter takes care of resetting a BMC
type BootDeviceSetter interface {
	BootDeviceSet(ctx context.Context, device string) (result string, err error)
}

// BMC management methods
type BMC interface {
	// NetworkSource() (result string, err repository.Error)
	CreateUser(context.Context) error
	UpdateUser(context.Context) error
	DeleteUser(context.Context) error
}

// BMCResetter options
type BMCResetter interface {
	// BMCReset resets the management console without rebooting the BMC (warm) or
	// Reboots the BMC (cold)
	BMCReset(ctx context.Context, rType string) error
}

// SetPower interface function for power actions
func SetPower(ctx context.Context, action string, m []PowerSetter) (result string, err error) {
	for _, elem := range m {
		result, setErr := elem.PowerSet(ctx, action)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return result, err
	}
	return result, multierror.Append(err, errors.New("power state failed"))
}

// SetBootDevice interface function for setting next boot device
func SetBootDevice(ctx context.Context, device string, m []BootDeviceSetter) (result string, err error) {
	for _, elem := range m {
		result, setErr := elem.BootDeviceSet(ctx, device)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return result, err
	}
	return result, multierror.Append(err, errors.New("set boot device failed"))
}

// CreateUser interface function
func CreateUser(ctx context.Context, u []BMC) (err error) {
	for _, elem := range u {
		setErr := elem.CreateUser(ctx)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return err
	}
	return multierror.Append(err, errors.New("create user failed"))
}

// UpdateUser interface function
func UpdateUser(ctx context.Context, u []BMC) (err error) {
	for _, elem := range u {
		setErr := elem.UpdateUser(ctx)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return err
	}
	return multierror.Append(err, errors.New("update user failed"))
}

// DeleteUser interface function
func DeleteUser(ctx context.Context, u []BMC) (err error) {
	for _, elem := range u {
		setErr := elem.DeleteUser(ctx)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return err
	}
	return multierror.Append(err, errors.New("delete user failed"))
}

// ResetBMC interface function
func ResetBMC(ctx context.Context, rType string, r []BMCResetter) (err error) {
	for _, elem := range r {
		setErr := elem.BMCReset(ctx, rType)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return err
	}
	return multierror.Append(err, errors.New("BMC reset failed"))
}
