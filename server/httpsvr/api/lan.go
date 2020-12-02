// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/bmc"
)

// updateLANConfig is the handler for the POST /ipmi-lan endpoint.
func updateLANConfig(c *gin.Context) {
	var req struct {
		IPSource bmc.IPSource `json:"ip_source" binding:"required"`
	}
	if c.BindJSON(&req) != nil {
		return
	}

	driver := bmc.NewDriverFromGinContext(c)
	if driver == nil {
		return
	}
	defer func() { _ = driver.Close() }()

	if err := driver.SetIPSource(req.IPSource); err != nil {
		c.Error(err) // nolint
		internalServerError(c)
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
