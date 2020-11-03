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

func (r *redfishBMC) Connect(ctx context.Context) repository.Error {
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
		return errMsg //nolint
	}
	r.conn = c
	return errMsg
}

func (r *redfishBMC) Close(ctx context.Context) {
	r.conn.Logout()
}

func (r *redfishBMC) on(ctx context.Context) (result string, errMsg repository.Error) {
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	for _, system := range ss {
		if system.PowerState == redfish.OnPowerState {
			break
		}
		err = system.Reset(redfish.OnResetType)
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			return "", errMsg
		}
	}
	return "on", errMsg
}

func (r *redfishBMC) off(ctx context.Context) (result string, errMsg repository.Error) {
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	for _, system := range ss {
		if system.PowerState == redfish.OffPowerState {
			break
		}
		err = system.Reset(redfish.GracefulShutdownResetType)
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			return "", errMsg
		}
	}
	return "off", errMsg
}

func (r *redfishBMC) status(ctx context.Context) (result string, errMsg repository.Error) {
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	for _, system := range ss {
		return string(system.PowerState), errMsg
	}
	return result, errMsg
}

func (r *redfishBMC) reset(ctx context.Context) (result string, errMsg repository.Error) {
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	for _, system := range ss {
		err = system.Reset(redfish.PowerCycleResetType)
		if err != nil {
			r.log.V(1).Info("warning", "msg", err.Error())
			r.off(ctx)
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
	return "reset", errMsg
}

func (r *redfishBMC) hardoff(ctx context.Context) (result string, errMsg repository.Error) {
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	for _, system := range ss {
		if system.PowerState == redfish.OnPowerState {
			break
		}
		err = system.Reset(redfish.ForceOffResetType)
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			return "", errMsg
		}
	}
	return "hardoff", errMsg
}

func (r *redfishBMC) cycle(ctx context.Context) (result string, errMsg repository.Error) {
	service := r.conn.Service
	ss, err := service.Systems()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	for _, system := range ss {
		err = system.Reset(redfish.GracefulRestartResetType)
		if err != nil {
			r.log.V(1).Info("warning", "msg", err.Error())
			r.off(ctx)
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
	return "cycle", errMsg
}
