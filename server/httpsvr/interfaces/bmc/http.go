// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package bmc

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tinkerbell/pbnj/server/httpsvr/util"
)

// NewDriverFromGinContext creates a new Driver using info in the http request
func NewDriverFromGinContext(c *gin.Context) Driver {
	driverType := util.FindDriver(c)
	driverOpts := DriverOptions{
		Address:  c.Param("ip"),
		Username: c.Request.Header.Get("X-IPMI-Username"),
		Password: c.Request.Header.Get("X-IPMI-Password"),
		Cipher:   -1,
	}
	if c.Request.Header.Get("X-IPMI-Cipher") != "" {
		cipher, err := strconv.Atoi(c.Request.Header.Get("X-IPMI-Cipher"))
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
		}
		driverOpts.Cipher = cipher
	}

	driver, err := NewDriver(c, driverType, driverOpts)
	if driver == nil {
		_ = c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return driver
}
