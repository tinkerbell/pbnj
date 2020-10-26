package machine

import (
	"context"
	"strconv"
	"time"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
	goipmi "github.com/vmware/goipmi"
)

type ipmiBootDevice struct {
	mAction  Action
	conn     *goipmi.Client
	user     string
	password string
	host     string
	port     string
	iface    string
}

func (b *ipmiBootDevice) Connect(ctx context.Context) repository.Error {
	result := make(chan repository.Error, 1)
	var errMsg repository.Error

	go func() {
		var port int
		port, err := strconv.Atoi(b.port)
		if err != nil {
			port = 623
		}
		var iface string
		if b.iface == "" {
			iface = "lanplus"
		} else {
			iface = b.iface
		}
		c := &goipmi.Connection{
			Hostname:  b.host,
			Username:  b.user,
			Password:  b.password,
			Port:      port,
			Interface: iface,
		}

		client, err := goipmi.NewClient(c)
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			result <- errMsg
			return
		}
		err = client.Open()
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			result <- errMsg
			return
		}
		b.conn = client
		result <- errMsg
		return
	}()

	// I can't figure out why this is needed, some kind of race condition possibly.
	// without it, the ctx.Done case will always trigger
	time.Sleep(100 * time.Millisecond)

	select {
	case <-ctx.Done():
		errMsg.Message = ctx.Err().Error()
		return errMsg
	case r := <-result:
		return r
	}
}

func (b *ipmiBootDevice) Close() {
	b.conn.Close()
}

func (b *ipmiBootDevice) setBootDevice() (result string, errMsg repository.Error) {
	var dev goipmi.BootDevice
	switch b.mAction.BootDeviceRequest.Device {
	case v1.DeviceRequest_NONE:
		dev = goipmi.BootDeviceNone
	case v1.DeviceRequest_BIOS:
		dev = goipmi.BootDeviceBios
	case v1.DeviceRequest_CDROM:
		dev = goipmi.BootDeviceCdrom
	case v1.DeviceRequest_DISK:
		dev = goipmi.BootDeviceDisk
	case v1.DeviceRequest_PXE:
		dev = goipmi.BootDevicePxe
	default:
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = "unknown boot device"
		return result, errMsg
	}
	err := b.conn.SetBootDevice(dev)
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return "", errMsg
	}
	return "boot device set: " + dev.String(), errMsg
}
