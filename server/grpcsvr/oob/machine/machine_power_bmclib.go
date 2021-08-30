package machine

import (
	"context"

	"github.com/bmc-toolbox/bmclib/devices"
	"github.com/bmc-toolbox/bmclib/discover"
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type bmclibBMC struct {
	log      logr.Logger
	conn     devices.Bmc
	user     string
	password string
	host     string
}

func (b *bmclibBMC) Connect(ctx context.Context) error {
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
	return b.conn.CheckCredentials()
}

func (b *bmclibBMC) Close(_ context.Context) {
	b.conn.Close()
}

func (b *bmclibBMC) PowerSet(ctx context.Context, action string) (result string, err error) {
	return doBmclibAction(ctx, action, b)
}

func doBmclibAction(ctx context.Context, action string, pwr *bmclibBMC) (result string, err error) {
	switch action {
	case v1.PowerAction_POWER_ACTION_ON.String():
		result, err = pwr.on(ctx)
	case v1.PowerAction_POWER_ACTION_OFF.String():
		result, err = pwr.off(ctx)
	case v1.PowerAction_POWER_ACTION_STATUS.String():
		result, err = pwr.status(ctx)
	case v1.PowerAction_POWER_ACTION_RESET.String():
		result, err = pwr.reset(ctx)
	case v1.PowerAction_POWER_ACTION_HARDOFF.String():
		result, err = pwr.hardoff(ctx)
	case v1.PowerAction_POWER_ACTION_CYCLE.String():
		result, err = pwr.cycle(ctx)
	case v1.PowerAction_POWER_ACTION_UNSPECIFIED.String():
		return result, &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "UNSPECIFIED power action",
		}
	default:
		return result, &repository.Error{
			Code:    v1.Code_value["INVALID_ARGUMENT"],
			Message: "unknown power action",
		}
	}
	return result, err
}

func (b *bmclibBMC) on(_ context.Context) (result string, err error) {
	ok, err := b.conn.PowerOn()
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	if !ok {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "error powering on",
		}
	}
	return On, nil
}

func (b *bmclibBMC) off(_ context.Context) (result string, err error) {
	ok, err := b.conn.PowerOff()
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	if !ok {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "error powering off",
		}
	}
	return Off, nil
}

func (b *bmclibBMC) status(_ context.Context) (result string, err error) {
	result, err = b.conn.PowerState()
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return result, nil
}

func (b *bmclibBMC) reset(_ context.Context) (result string, err error) {
	ok, err := b.conn.PowerCycle()
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	if !ok {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "error with power reset",
		}
	}
	return Reset, nil
}

func (b *bmclibBMC) hardoff(_ context.Context) (result string, err error) {
	ok, err := b.conn.PowerOff()
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	if !ok {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "error with power hardoff",
		}
	}
	return "hardoff", nil
}

func (b *bmclibBMC) cycle(_ context.Context) (result string, err error) {
	ok, err := b.conn.PowerCycle()
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	if !ok {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: "error with power cycle",
		}
	}
	return "cycle", nil
}
