package machine

import (
	"context"
	"fmt"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type bmclibClient struct {
	log      logr.Logger
	conn     *bmclib.Client
	user     string
	password string
	host     string
}

func (b *bmclibClient) Connect(ctx context.Context) error {
	b.conn = bmclib.NewClient(b.host, "623", b.user, b.password, bmclib.WithLogger(b.log))
	//b.conn.Registry.Drivers = b.conn.Registry.FilterForCompatible(ctx)
	return b.conn.Open(ctx)
}

func (b *bmclibClient) Close(ctx context.Context) {
	b.conn.Close(ctx)
}

func (b *bmclibClient) PowerSet(ctx context.Context, action string) (result string, err error) {
	var pwrAction string
	switch action {
	case v1.PowerAction_POWER_ACTION_ON.String():
		pwrAction = "on"
	case v1.PowerAction_POWER_ACTION_OFF.String():
		pwrAction = "off"
	case v1.PowerAction_POWER_ACTION_STATUS.String():
		pwrAction = "status"
	case v1.PowerAction_POWER_ACTION_RESET.String():
		pwrAction = "reset"
	case v1.PowerAction_POWER_ACTION_HARDOFF.String():
		pwrAction = "off"
	case v1.PowerAction_POWER_ACTION_CYCLE.String():
		pwrAction = "cycle"
	case v1.PowerAction_POWER_ACTION_UNSPECIFIED.String():
		return "", &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "UNSPECIFIED power action",
		}
	default:
		return "", &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: fmt.Sprintf("unknown power action: %v", action),
		}
	}
	var ok bool
	if pwrAction == "status" {
		pwrAction, err = b.conn.GetPowerState(ctx)
	} else {
		ok, err = b.conn.SetPowerState(ctx, pwrAction)
	}

	if err != nil {
		return "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	if !ok {
		return "", &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: fmt.Sprintf("error completing power %v action", pwrAction),
		}
	}
	return pwrAction, nil
}
