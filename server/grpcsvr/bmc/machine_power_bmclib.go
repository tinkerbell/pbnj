package bmc

import (
	"github.com/bmc-toolbox/bmclib/devices"
	"github.com/bmc-toolbox/bmclib/discover"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
)

type bmclibBMC struct {
	mAction  MachineAction
	conn     devices.Bmc
	user     string
	password string
	host     string
}

func (b *bmclibBMC) connection() repository.Error {
	var errMsg repository.Error
	l := b.mAction.Log.GetContextLogger(b.mAction.Ctx)
	connection, err := discover.ScanAndConnect(b.host, b.user, b.password, discover.WithLogger(l))
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return errMsg //nolint
	}
	switch conn := connection.(type) {
	case devices.Bmc:
		b.conn = conn
	default:
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "Unknown device"
		return errMsg //nolint
	}
	return errMsg //nolint
}

func (b *bmclibBMC) close() {
	b.conn.Close()
}

func (b *bmclibBMC) on() (result string, errMsg repository.Error) {
	ok, err := b.conn.PowerOn()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	if !ok {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "error powering on"
		return "", errMsg
	}
	return "on", errMsg
}

func (b *bmclibBMC) off() (result string, errMsg repository.Error) {
	ok, err := b.conn.PowerOff()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	if !ok {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "error powering off"
		return "", errMsg
	}
	return "off", errMsg
}

func (b *bmclibBMC) status() (result string, errMsg repository.Error) {
	result, err := b.conn.PowerState()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return result, errMsg
	}
	return result, errMsg
}

func (b *bmclibBMC) reset() (result string, errMsg repository.Error) {
	ok, err := b.conn.PowerCycle()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	if !ok {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "error with power reset"
		return "", errMsg
	}
	return "reset", errMsg
}

func (b *bmclibBMC) hardoff() (result string, errMsg repository.Error) {
	ok, err := b.conn.PowerOff()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	if !ok {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "error with power hardoff"
		return "", errMsg
	}
	return "hardoff", errMsg
}

func (b *bmclibBMC) cycle() (result string, errMsg repository.Error) {
	ok, err := b.conn.PowerCycle()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	if !ok {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "error with power cycle"
		return "", errMsg
	}
	return "cycle", errMsg
}
