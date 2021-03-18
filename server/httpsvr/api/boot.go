// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/boot"
)

// updateBootOptions is the handler for the PATCH /boot endpoint.
func updateBootOptions(c *gin.Context) {
	var opts boot.Options
	if c.BindJSON(&opts) != nil {
		return
	}

	driver := boot.NewDriverFromGinContext(c)
	if driver == nil {
		return
	}
	defer func() { _ = driver.Close() }()

	if err := driver.SetBootOptions(opts); err != nil {
		c.Error(err) // nolint
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
