package oob

import (
	"context"
	"errors"

	"github.com/hashicorp/go-multierror"
)

// Machine management methods
type Machine interface {
	// BootDevice sets the next boot device
	BootDevice(ctx context.Context, device string) (result string, err error)
	// Power get status and sets power states like on/off/etc
	Power(ctx context.Context, action string) (result string, err error)
}

// BMC management methods
type BMC interface {
	// NetworkSource() (result string, err repository.Error)
	CreateUser(context.Context) error
	UpdateUser(context.Context) error
	DeleteUser(context.Context) error
}

// BMCReset options
type BMCReset interface {
	// ResetWarm resets the management console without rebooting the BMC
	ResetWarm(context.Context) error
	// ResetCold Reboots the BMC
	ResetCold(context.Context) error
}

// MachinePower interface function for power actions
func MachinePower(ctx context.Context, action string, m []Machine) (result string, err error) {
	for _, elem := range m {
		result, setErr := elem.Power(ctx, action)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return result, err
	}
	return result, multierror.Append(err, errors.New("power state failed"))
}

// MachineBootDevice interface function for setting next boot device
func MachineBootDevice(ctx context.Context, device string, m []Machine) (result string, err error) {
	for _, elem := range m {
		result, setErr := elem.BootDevice(ctx, device)
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

// WarmBMCReset interface function
func WarmBMCReset(ctx context.Context, r []BMCReset) (err error) {
	for _, elem := range r {
		setErr := elem.ResetWarm(ctx)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return err
	}
	return multierror.Append(err, errors.New("BMC warm reset failed"))
}

// ColdBMCReset interface function
func ColdBMCReset(ctx context.Context, r []BMCReset) (err error) {
	for _, elem := range r {
		setErr := elem.ResetCold(ctx)
		if setErr != nil {
			err = multierror.Append(err, setErr)
			continue
		}
		return err
	}
	return multierror.Append(err, errors.New("BMC cold reset failed"))
}
