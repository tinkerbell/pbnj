package oob

import (
	"context"
	"errors"

	"github.com/hashicorp/go-multierror"
)

// BMC management methods
type BMC interface {
	// NetworkSource() (result string, err repository.Error)
	CreateUser(context.Context) error
	UpdateUser(context.Context) error
	DeleteUser(context.Context) error
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
