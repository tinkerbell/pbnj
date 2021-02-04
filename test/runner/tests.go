package runner

import (
	"time"

	v1 "github.com/tinkerbell/pbnj/api/v1"
)

var (
	powerRequestON     = v1.PowerAction_POWER_ACTION_ON
	powerRequestOFF    = v1.PowerAction_POWER_ACTION_OFF
	powerRequestSTATUS = v1.PowerAction_POWER_ACTION_STATUS
	powerRequestCYCLE  = v1.PowerAction_POWER_ACTION_CYCLE
	//powerRequestRESET   = v1.PowerAction_POWER_ACTION_RESET
	//powerRequestHARDOFF = v1.PowerAction_POWER_ACTION_HARDOFF
	//deviceRequestNONE   = v1.BootDevice_BOOT_DEVICE_NONE
	deviceRequestBIOS = v1.BootDevice_BOOT_DEVICE_BIOS
	//deviceRequestDISK   = v1.BootDevice_BOOT_DEVICE_DISK
	//deviceRequestCDROM  = v1.BootDevice_BOOT_DEVICE_CDROM
	deviceRequestPXE = v1.BootDevice_BOOT_DEVICE_PXE
	lookup           = map[string]map[string]expected{
		"happyTests":    happyTests,
		"userMgmtTests": userMgmtTests,
		//"notIdentifiableTests": notIdentifiableTests,
	}
	happyTests = map[string]expected{
		"1 power off": {
			Action: powerRequestOFF,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "power action: POWER_ACTION_OFF",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "off",
				Complete:    true,
				Messages:    []string{"working on power POWER_ACTION_OFF", "connected to BMC", "power POWER_ACTION_OFF complete"},
			},
			WaitTime: 15 * time.Second,
		},
		"2 power status off": {
			Action: powerRequestSTATUS,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "power action: POWER_ACTION_STATUS",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "off - soft",
				Complete:    true,
				Messages:    []string{"working on power POWER_ACTION_STATUS", "connected to BMC", "power POWER_ACTION_STATUS complete"},
			},
		},
		"3 set device bios": {
			Action: deviceRequestBIOS,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "setting boot device",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "complete",
				Complete:    true,
				Messages:    []string{"working on setting boot device: BOOT_DEVICE_BIOS", "connecting to BMC", "setting boot device: BOOT_DEVICE_BIOS complete"},
			},
			WaitTime: 1 * time.Second,
		},
		"4 power on": {
			Action: powerRequestON,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "power action: POWER_ACTION_ON",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "on",
				Complete:    true,
				Messages:    []string{"working on power POWER_ACTION_ON", "connected to BMC", "power POWER_ACTION_ON complete"},
			},
			WaitTime: 1 * time.Second,
		},
		"5 power status on": {
			Action: powerRequestSTATUS,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "power action: POWER_ACTION_STATUS",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "on",
				Complete:    true,
				Messages:    []string{"working on power POWER_ACTION_STATUS", "connected to BMC", "power POWER_ACTION_STATUS complete"},
			},
			WaitTime: 1 * time.Second,
		},
		"6 set device pxe": {
			Action: deviceRequestPXE,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "setting boot device",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "complete",
				Complete:    true,
				Messages:    []string{"working on setting boot device: BOOT_DEVICE_PXE", "connecting to BMC", "setting boot device: BOOT_DEVICE_PXE complete"},
			},
			WaitTime: 1 * time.Second,
		},
		"7 power cycle": {
			Action: powerRequestCYCLE,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "power action: POWER_ACTION_CYCLE",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "cycle",
				Complete:    true,
				Messages:    []string{"working on power POWER_ACTION_CYCLE", "connected to BMC", "power POWER_ACTION_CYCLE complete"},
			},
			WaitTime: 60 * time.Second,
		},
		"8 power status on": {
			Action: powerRequestSTATUS,
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "power action: POWER_ACTION_STATUS",
				Error:       &v1.Error{},
				State:       "complete",
				Result:      "on",
				Complete:    true,
				Messages:    []string{"working on power POWER_ACTION_STATUS", "connected to BMC", "power POWER_ACTION_STATUS complete"},
			},
		},

		//"power hardoff": {Action: &PowerRequest_HARDOFF, Want: notImplementedWant("HARD OFF")},
		//"power reset":   {Action: &PowerRequest_RESET, Want: notImplementedWant("RESET")},
	}
	userMgmtTests = map[string]expected{
		"1 create a user": {
			Action: &v1.CreateUserRequest{
				UserCreds: &v1.UserCreds{
					Username: "jacob",
					Password: "Jacob1234",
					UserRole: v1.UserRole_USER_ROLE_ADMIN,
				},
			},
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "creating user",
				Error:       &v1.Error{},
				State:       "complete",
				Complete:    true,
				Messages:    []string{"working on creating user: jacob", "connecting to BMC", "connected to BMC", "creating user: jacob complete"},
			},
		},
		"2 delete a user": {
			Action: &v1.DeleteUserRequest{
				Username: "jacob",
			},
			Want: &v1.StatusResponse{
				Id:          "12345",
				Description: "deleting user",
				Error:       &v1.Error{},
				State:       "complete",
				Complete:    true,
				Messages:    []string{"working on deleting user: jacob", "connecting to BMC", "connected to BMC", "deleting user: jacob complete"},
			},
		},
	}
	/*
			deviceHappyTests = map[string]expected{
				"1 set device pxe": {
					Action: deviceRequestPXE,
					Want: &v1.StatusResponse{
						Id:          "12345",
						Description: "power action: OFF",
						Error:       &v1.Error{},
						State:       "complete",
						Result:      "off",
						Complete:    true,
						Messages:    []string{"working on power OFF", "connecting to BMC", "connected to BMC", "power OFF complete"},
					},
					WaitTime: 1 * time.Second,
				},
				"2 power status": {
					Action: powerRequestSTATUS,
					Want: &v1.StatusResponse{
						Id:          "12345",
						Description: "power action: STATUS",
						Error:       &v1.Error{},
						State:       "complete",
						Result:      "off",
						Complete:    true,
						Messages:    []string{"working on power STATUS", "connecting to BMC", "connected to BMC", "power STATUS complete"},
					},
				},
				"3 power on": {
					Action: powerRequestON,
					Want: &v1.StatusResponse{
						Id:          "12345",
						Description: "power action: ON",
						Error:       &v1.Error{},
						State:       "complete",
						Result:      "on",
						Complete:    true,
						Messages:    []string{"working on power ON", "connecting to BMC", "connected to BMC", "power ON complete"},
					},
					WaitTime: 180 * time.Second,
				},
				"4 power status": {
					Action: powerRequestSTATUS,
					Want: &v1.StatusResponse{
						Id:          "12345",
						Description: "power action: STATUS",
						Error:       &v1.Error{},
						State:       "complete",
						Result:      "on",
						Complete:    true,
						Messages:    []string{"working on power STATUS", "connecting to BMC", "connected to BMC", "power STATUS complete"},
					},
				},
				"5 power cycle": {
					Action: powerRequestCYCLE,
					Want: &v1.StatusResponse{
						Id:          "12345",
						Description: "power action: CYCLE",
						Error:       &v1.Error{},
						State:       "complete",
						Result:      "cycle",
						Complete:    true,
						Messages:    []string{"working on power CYCLE", "connecting to BMC", "connected to BMC", "power CYCLE complete"},
					},
					WaitTime: 60 * time.Second,
				},
				"6 power status": {
					Action: powerRequestSTATUS,
					Want: &v1.StatusResponse{
						Id:          "12345",
						Description: "power action: STATUS",
						Error:       &v1.Error{},
						State:       "complete",
						Result:      "on",
						Complete:    true,
						Messages:    []string{"working on power STATUS", "connecting to BMC", "connected to BMC", "power STATUS complete"},
					},
				},
			}
		notIdentifiableTests = map[string]expected{
			"power status":  {Action: powerRequestSTATUS, Want: notIdentifiableWant},
			"power on":      {Action: powerRequestON, Want: notIdentifiableWant},
			"power off":     {Action: powerRequestOFF, Want: notIdentifiableWant},
			"power hardoff": {Action: powerRequestHARDOFF, Want: notIdentifiableWant},
			"power cycle":   {Action: powerRequestCYCLE, Want: notIdentifiableWant},
			"power reset":   {Action: powerRequestRESET, Want: notIdentifiableWant},
		}
		notIdentifiableWant = &v1.StatusResponse{
			Id:          "12345",
			Description: "power action",
			Error: &v1.Error{
				Code:    2,
				Message: "unable to identify the vendor",
				Details: nil,
			},
			State:    "complete",
			Result:   "action failed",
			Complete: true,
			Messages: []string{"connecting to BMC", "connecting to BMC failed"},
		}
	*/
)
