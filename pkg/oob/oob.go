package oob

import (
	"context"
	"errors"

	"github.com/hashicorp/go-multierror"
)

// BootDeviceSetter takes care of resetting a BMC
type BootDeviceSetter interface {
	BootDeviceSet(ctx context.Context, device string, persistent, efiBoot bool) (result string, err error)
}

// BMC management methods
type BMC interface {
	// NetworkSource() (result string, err repository.Error)
	CreateUser(context.Context) error
	UpdateUser(context.Context) error
	DeleteUser(context.Context) error
}

// SetBootDevice interface function for setting next boot device
func SetBootDevice(ctx context.Context, device string, persistent, efiBoot bool, m []BootDeviceSetter) (result string, err error) {
	for _, elem := range m {
		select {
		case <-ctx.Done():
			err = multierror.Append(err, ctx.Err())
			break
		default:
			if elem != nil {
				result, setErr := elem.BootDeviceSet(ctx, device, persistent, efiBoot)
				if setErr != nil {
					err = multierror.Append(err, setErr)
					continue
				}
				return result, nil
			}
			err = multierror.Append(err, errors.New("set boot device request not executed"))
		}
	}
	return result, multierror.Append(err, errors.New("set boot device failed"))
}

// CreateUser interface function
func CreateUser(ctx context.Context, u []BMC) (err error) {
	for _, elem := range u {
		select {
		case <-ctx.Done():
			err = multierror.Append(err, ctx.Err())
			break
		default:
			if elem != nil {
				setErr := elem.CreateUser(ctx)
				if setErr != nil {
					err = multierror.Append(err, setErr)
					continue
				}
				return nil
			}
			err = multierror.Append(err, errors.New("create user request not executed"))
		}
	}
	return multierror.Append(err, errors.New("create user failed"))
}

// UpdateUser interface function
func UpdateUser(ctx context.Context, u []BMC) (err error) {
	for _, elem := range u {
		select {
		case <-ctx.Done():
			err = multierror.Append(err, ctx.Err())
			break
		default:
			if elem != nil {
				setErr := elem.UpdateUser(ctx)
				if setErr != nil {
					err = multierror.Append(err, setErr)
					continue
				}
				return nil
			}
			err = multierror.Append(err, errors.New("update user request not executed"))
		}
	}
	return multierror.Append(err, errors.New("update user failed"))
}

// DeleteUser interface function
func DeleteUser(ctx context.Context, u []BMC) (err error) {
	for _, elem := range u {
		select {
		case <-ctx.Done():
			err = multierror.Append(err, ctx.Err())
			break
		default:
			if elem != nil {
				setErr := elem.DeleteUser(ctx)
				if setErr != nil {
					err = multierror.Append(err, setErr)
					continue
				}
				return nil
			}
			err = multierror.Append(err, errors.New("delete user request not executed"))
		}
	}
	return multierror.Append(err, errors.New("delete user failed"))
}
