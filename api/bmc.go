// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinkerbell/pbnj/interfaces/bmc"
)

// bmcAction is the handler for the POST /bmc endpoint.
func bmcAction(c *gin.Context) {
	var req bmc.BmcRequest
	if c.BindJSON(&req) != nil {
		return
	}

	driver := bmc.NewDriverFromGinContext(c)
	if driver == nil {
		return
	}
	defer func() { _ = driver.Close() }()

	if err := driver.BMC(req); err != nil {
		c.Error(err)
		internalServerError(c)
		return
	}
	c.Writer.WriteHeader(http.StatusNoContent)
}
