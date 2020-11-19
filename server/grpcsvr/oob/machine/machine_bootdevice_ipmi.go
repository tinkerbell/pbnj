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
func (b *ipmiBootDevice) Connect(ctx context.Context) error {
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
		return &errMsg
	}
	err = client.Open()
	if err != nil {
		errMsg.Code = v1.Code_value["UNKNOWN"]
		errMsg.Message = err.Error()
		return &errMsg
	}
	b.conn = client
	return nil
}

// Close the connection to a BMC
func (b *ipmiBootDevice) Close(ctx context.Context) {
	b.conn.Close()
}

// setBootDevice will try to set boot device using ipmitool interface lan and lanplus
func (b *ipmiBootDevice) BootDevice(ctx context.Context, device string) (result string, err error) {
	var errMsg repository.Error
	var dev goipmi.BootDevice
	switch device {
	case v1.BootDevice_BOOT_DEVICE_NONE.String():
		dev = goipmi.BootDeviceNone
	case v1.BootDevice_BOOT_DEVICE_BIOS.String():
		dev = goipmi.BootDeviceBios
	case v1.BootDevice_BOOT_DEVICE_CDROM.String():
		dev = goipmi.BootDeviceCdrom
	case v1.BootDevice_BOOT_DEVICE_DISK.String():
		dev = goipmi.BootDeviceDisk
	case v1.BootDevice_BOOT_DEVICE_PXE.String():
		dev = goipmi.BootDevicePxe
	case v1.BootDevice_BOOT_DEVICE_UNSPECIFIED.String():
		errMsg.Code = v1.Code_value["INVALID_ARGUMENT"]
		errMsg.Message = "UNSPECIFIED boot device"
		return result, &errMsg
	default:
		errMsg.Code = v1.Code_value["INVALID_ARGUMENT"]
		errMsg.Message = "unknown boot device"
		return result, &errMsg
	}
	setErr := b.conn.SetBootDevice(dev)
	if setErr != nil {
		var iface string
		if b.iface == lan {
			iface = lanplus
		} else if b.iface == lanplus {
			iface = lan
		}
		b.iface = iface
		connErr := b.Connect(ctx)
		if connErr != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = connErr.Error()
			return result, &errMsg
		}
		conn2err := b.conn.SetBootDevice(dev)
		if conn2err != nil {
			errMsg.Code = v1.Code_value["UNKNOWN"]
			errMsg.Message = conn2err.Error()
			return result, &errMsg
		}
	}

	return "boot device set: " + dev.String(), nil
}
