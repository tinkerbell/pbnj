package machine

import (
	"context"

	"github.com/gebn/bmc"
	"github.com/gebn/bmc/pkg/ipmi"
	"github.com/go-logr/logr"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type ipmiBMC struct {
	log       logr.Logger
	transport bmc.SessionlessTransport
	conn      bmc.Session
	user      string
	password  string
	host      string
}

func (b *ipmiBMC) Connect(ctx context.Context) error {
	var errMsg repository.Error
	machine, err := bmc.Dial(ctx, b.host)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return &errMsg
	}
	b.transport = machine

	sess, err := machine.NewSession(ctx, &bmc.SessionOpts{
		Username:          b.user,
		Password:          []byte(b.password),
		MaxPrivilegeLevel: ipmi.PrivilegeLevelOperator,
	})
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return &errMsg
	}
	b.conn = sess
	return nil
}

func (b *ipmiBMC) Close(ctx context.Context) {
	b.transport.Close()
	b.conn.Close(ctx)
}

func (b *ipmiBMC) Power(ctx context.Context, action string) (result string, err error) {
	return doIpmiAction(ctx, action, b)
}

func doIpmiAction(ctx context.Context, action string, pwr *ipmiBMC) (result string, err error) {
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

func (b *ipmiBMC) on(ctx context.Context) (result string, err error) {
	err = b.conn.ChassisControl(ctx, ipmi.ChassisControlPowerOn)
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return "on", nil
}

func (b *ipmiBMC) off(ctx context.Context) (result string, err error) {
	err = b.conn.ChassisControl(ctx, ipmi.ChassisControlSoftPowerOff)
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return "off", nil
}

func (b *ipmiBMC) status(ctx context.Context) (result string, err error) {
	result = "off"
	status, err := b.conn.GetChassisStatus(ctx)
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	if status.PoweredOn {
		result = "on"
	}
	return result, nil
}

func (b *ipmiBMC) reset(ctx context.Context) (result string, err error) {
	err = b.conn.ChassisControl(ctx, ipmi.ChassisControlHardReset)
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return "reset", nil
}

func (b *ipmiBMC) hardoff(ctx context.Context) (result string, err error) {
	err = b.conn.ChassisControl(ctx, ipmi.ChassisControlPowerOff)
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return "hardoff", nil
}

func (b *ipmiBMC) cycle(ctx context.Context) (result string, err error) {
	err = b.conn.ChassisControl(ctx, ipmi.ChassisControlPowerCycle)
	if err != nil {
		return result, &repository.Error{
			Code:    v1.Code_value["UNKNOWN"],
			Message: err.Error(),
		}
	}
	return "cycle", nil
}
