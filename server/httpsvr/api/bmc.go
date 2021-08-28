// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/bmc"
)

// bmcAction is the handler for the POST /bmc endpoint.
func bmcAction(c *gin.Context) {
	var req struct {
		Action bmc.Action `json:"action" binding:"required"`
	}
	if c.BindJSON(&req) != nil {
		return
	}

	driver := bmc.NewDriverFromGinContext(c)
	if driver == nil {
		return
	}
	defer func() { _ = driver.Close() }()

	if err := driver.BMC(req.Action); err != nil {
		internalServerError(c, err)
		return
	}
	c.Writer.WriteHeader(http.StatusNoContent)
}
