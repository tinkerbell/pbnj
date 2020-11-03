package machine

import (
	"context"
	"net"
	"strconv"
	"strings"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/repository"
	goipmi "github.com/vmware/goipmi"
)

const (
	lan     = "lan"
	lanplus = "lanplus"
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

// Connect to BMC using ipmitool
func (b *ipmiBootDevice) Connect(ctx context.Context) repository.Error {
	var errMsg repository.Error

	if strings.Contains(b.host, ":") {
		host, port, err := net.SplitHostPort(b.host)
		if err == nil {
			b.host = host
			b.port = port
		}
	}
	var port int
	port, err := strconv.Atoi(b.port)
	if err != nil {
		port = 623
	}
	var iface string
	if b.iface == "" {
		iface = lanplus
		b.iface = lanplus
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
		return errMsg
	}
	err = client.Open()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return errMsg
	}
	b.conn = client
	return errMsg
}

// Close the connection to a BMC
func (b *ipmiBootDevice) Close(ctx context.Context) {
	b.conn.Close()
}

// setBootDevice will try to set boot device using ipmitool interface lan and lanplus
func (b *ipmiBootDevice) setBootDevice(ctx context.Context) (result string, errMsg repository.Error) {
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
		var iface string
		if b.iface == lan {
			iface = lanplus
		} else if b.iface == lanplus {
			iface = lan
		}
		b.iface = iface
		errMsg = b.Connect(ctx)
		if errMsg.Message != "" {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			return "", errMsg
		}
		err := b.conn.SetBootDevice(dev)
		if err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = err.Error()
			return "", errMsg
		}
	}

	return "boot device set: " + dev.String(), errMsg
}
