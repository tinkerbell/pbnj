// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"github.com/gin-gonic/gin"
)

// FindDriver finds the driver type given a manufacturer, defaulting to ipmitool
func FindDriver(c *gin.Context) string {
	manufacturer := c.Request.Header.Get("X-DEVICE-MANUFACTURER")
	if manufacturer == "dell" {
		return "racadm"
	}
	return "ipmitool"
}
