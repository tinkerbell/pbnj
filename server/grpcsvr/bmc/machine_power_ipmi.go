package bmc

import (
	"github.com/gebn/bmc"
	"github.com/gebn/bmc/pkg/ipmi"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type ipmiBMC struct {
	mAction   MachineAction
	transport bmc.SessionlessTransport
	conn      bmc.Session
	user      string
	password  string
	host      string
}

func (b *ipmiBMC) connection() repository.Error {
	var errMsg repository.Error
	machine, err := bmc.Dial(b.mAction.Ctx, b.host)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return errMsg
	}
	b.transport = machine

	sess, err := machine.NewSession(b.mAction.Ctx, &bmc.SessionOpts{
		Username:          b.user,
		Password:          []byte(b.password),
		MaxPrivilegeLevel: ipmi.PrivilegeLevelOperator,
	})
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return errMsg
	}
	b.conn = sess
	return errMsg
}

func (b *ipmiBMC) close() {
	b.transport.Close()
	b.conn.Close(b.mAction.Ctx)
}

func (b *ipmiBMC) on() (result string, errMsg repository.Error) {
	err := b.conn.ChassisControl(b.mAction.Ctx, ipmi.ChassisControlPowerOn)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	return "on", errMsg
}

func (b *ipmiBMC) off() (result string, errMsg repository.Error) {
	err := b.conn.ChassisControl(b.mAction.Ctx, ipmi.ChassisControlSoftPowerOff)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	return "off", errMsg
}

func (b *ipmiBMC) status() (result string, errMsg repository.Error) {
	result = "off"
	status, err := b.conn.GetChassisStatus(b.mAction.Ctx)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	if status.PoweredOn {
		result = "on"
	}
	return result, errMsg
}

func (b *ipmiBMC) reset() (result string, errMsg repository.Error) {
	err := b.conn.ChassisControl(b.mAction.Ctx, ipmi.ChassisControlHardReset)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	return "reset", errMsg
}

func (b *ipmiBMC) hardoff() (result string, errMsg repository.Error) {
	err := b.conn.ChassisControl(b.mAction.Ctx, ipmi.ChassisControlPowerOff)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	return "hardoff", errMsg
}

func (b *ipmiBMC) cycle() (result string, errMsg repository.Error) {
	err := b.conn.ChassisControl(b.mAction.Ctx, ipmi.ChassisControlPowerCycle)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	return "cycle", errMsg
}
