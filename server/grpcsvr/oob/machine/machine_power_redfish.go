package machine

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type redfishBMC struct {
	log      logr.Logger
	conn     *gofish.APIClient
	user     string
	password string
	host     string
}

func (r *redfishBMC) Connect(ctx context.Context) error {
	var errMsg repository.Error
	config := gofish.ClientConfig{
		Endpoint: "https://" + r.host,
		Username: r.user,
		Password: r.password,
		Insecure: true,
	}

	c, err := gofish.Connect(config)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return &errMsg
	}
	r.conn = c
	return nil
}

func (r *redfishBMC) Close(ctx context.Context) {
	r.conn.Logout()
}

func (r *redfishBMC) PowerSet(ctx context.Context, action string) (result string, err error) {
	return doRedfishAction(ctx, action, r)
}

func doRedfishAction(ctx context.Context, action string, pwr *redfishBMC) (result string, err error) {
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

func (r *redfishBMC) on(ctx context.Context) (result string, err error) {
	var errMsg repository.Error
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return result, &errMsg
	}
	for _, system := range ss {
		if system.PowerState == redfish.OnPowerState {
			break
		}
		err = system.Reset(redfish.OnResetType)
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			return result, &errMsg
		}
	}
	return "on", nil
}

func (r *redfishBMC) off(ctx context.Context) (result string, err error) {
	var errMsg repository.Error
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return result, &errMsg
	}
	for _, system := range ss {
		if system.PowerState == redfish.OffPowerState {
			break
		}
		err = system.Reset(redfish.GracefulShutdownResetType)
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			return result, &errMsg
		}
	}
	return "off", nil
}

func (r *redfishBMC) status(ctx context.Context) (result string, err error) {
	var errMsg repository.Error
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return result, &errMsg
	}
	for _, system := range ss {
		return string(system.PowerState), &errMsg
	}
	return result, nil
}

func (r *redfishBMC) reset(ctx context.Context) (result string, err error) {
	var errMsg repository.Error
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return result, &errMsg
	}
	for _, system := range ss {
		err = system.Reset(redfish.PowerCycleResetType)
		if err != nil {
			r.log.V(1).Info("warning", "msg", err.Error())
			_, _ = r.off(ctx)
			for wait := 1; wait < 10; wait++ {
				status, _ := r.status(ctx)
				if status == "off" {
					break
				}
				time.Sleep(1 * time.Second)
			}
			_, errMsg := r.on(ctx)
			return "reset", errMsg
		}
	}
	return "reset", nil
}

func (r *redfishBMC) hardoff(ctx context.Context) (result string, err error) {
	var errMsg repository.Error
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return result, &errMsg
	}
	for _, system := range ss {
		if system.PowerState == redfish.OnPowerState {
			break
		}
		err = system.Reset(redfish.ForceOffResetType)
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			return result, &errMsg
		}
	}
	return "hardoff", nil
}

func (r *redfishBMC) cycle(ctx context.Context) (result string, err error) {
	var errMsg repository.Error
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return result, &errMsg
	}
	for _, system := range ss {
		err = system.Reset(redfish.GracefulRestartResetType)
		if err != nil {
			r.log.V(1).Info("warning", "msg", err.Error())
			_, _ = r.off(ctx)
			for wait := 1; wait < 10; wait++ {
				status, _ := r.status(ctx)
				if status == "off" {
					break
				}
				time.Sleep(1 * time.Second)
			}
			_, errMsg := r.on(ctx)
			return "cycle", errMsg
		}
	}
	return "cycle", nil
}
