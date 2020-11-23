package bmc

import (
	"context"

	"github.com/bmc-toolbox/bmclib/devices"
	"github.com/bmc-toolbox/bmclib/discover"
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type bmcilbResetBMC struct {
	conn     devices.Bmc
	host     string
	log      logr.Logger
	password string
	user     string
}

func (b *bmcilbResetBMC) Connect(ctx context.Context) error {
	connection, err := discover.ScanAndConnect(b.host, b.user, b.password, discover.WithLogger(b.log), discover.WithContext(ctx))
	if err != nil {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	switch conn := connection.(type) {
	case devices.Bmc:
		b.conn = conn
	default:
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "Unknown device",
		}
	}
	return nil
}

func (b *bmcilbResetBMC) Close(ctx context.Context) {
	b.conn.Close()
}

func (b *bmcilbResetBMC) BMCReset(ctx context.Context, rType string) error {
	var err error
	var ok bool
	switch rType {
	case v1.ResetKind_RESET_KIND_COLD.String():
		var coldErr error
		ok, coldErr = b.conn.PowerCycleBmc()
		if coldErr != nil {
			err = &repository.Error{
				Code:    v1.Code_value["UNKNOWN"],
				Message: coldErr.Error(),
			}
		}
	case v1.ResetKind_RESET_KIND_WARM.String():
		err = &repository.Error{
			Code:    v1.Code_value["UNIMPLEMENTED"],
			Message: "reset warm unimplemented",
		}
	case v1.ResetKind_RESET_KIND_UNSPECIFIED.String():
		err = &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "UNSPECIFIED reset request",
		}
	default:
		err = &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "unknown reset request",
		}
	}

	if err != nil {
		return err
	}
	if !ok {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "reset failed",
		}
	}
	return nil
}
