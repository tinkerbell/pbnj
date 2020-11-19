package bmc

import (
	"context"
	"fmt"

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
	b.log.V(0).Info("debugging", "struct", fmt.Sprintf("%+v", b))
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
	b.log.V(0).Info("debugging", "struct", fmt.Sprintf("%+v", b))
	return nil
}

func (b *bmcilbResetBMC) Close(ctx context.Context) {
	b.conn.Close()
}

func (b *bmcilbResetBMC) ResetCold(ctx context.Context) error {
	ok, err := b.conn.PowerCycleBmc()
	if err != nil {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	if !ok {
		return &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "reset failed",
		}
	}
	return nil
}

func (b *bmcilbResetBMC) ResetWarm(ctx context.Context) error {
	return &repository.Error{
		Code:    v1.Code_value["UNIMPLEMENTED"],
		Message: "reset warm unimplemented",
	}

}
