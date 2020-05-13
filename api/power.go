// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tinkerbell/pbnj/interfaces/power"
	"github.com/tinkerbell/pbnj/reqid"
)

// powerAction is the handler for the POST /power endpoint.
func powerAction(c *gin.Context) {
	var req struct {
		Action      power.Operation `json:"action" binding:"required"`
		SoftTimeout string          `json:"soft_timeout,omitempty"`
		OffDuration string          `json:"off_duration,omitempty"`
	}
	if c.BindJSON(&req) != nil {
		return
	}

	opts := power.DefaultOptions

	if req.SoftTimeout != "" {
		d, err := time.ParseDuration(req.SoftTimeout)
		if err != nil {
			badRequest(c, err)
			return
		}
		opts.SoftTimeout = d
	}

	if req.OffDuration != "" {
		d, err := time.ParseDuration(req.OffDuration)
		if err != nil {
			badRequest(c, err)
			return
		}
		opts.OffDuration = d
	}

	if c.Request.Header.Get("X-DEVICE-MANUFACTURER") == "intel" {
		// Intel BMC returns error when attempting to turn off
		// server that is already off
		opts.IgnoreRunError = true

		// In some cases, Intel BMC silently ignores attempt to
		// turn on server if command is sent within a few seconds
		// of server being turned off
		minIntelOffDuration := 10 * time.Second
		if opts.OffDuration < minIntelOffDuration {
			opts.OffDuration = minIntelOffDuration
		}
	}

	driver := power.NewDriverFromGinContext(c)
	if driver == nil {
		return
	}

	id := reqid.FromContext(c)
	task := power.StartTask(c, id, req.Action, driver, opts)
	renderTaskStarted(c, task)
}

// powerStatus is the handler for the GET /power endpoint.
func powerStatus(c *gin.Context) {
	driver := power.NewDriverFromGinContext(c)
	if driver == nil {
		return
	}
	defer func() { _ = driver.Close() }()

	if status, err := driver.PowerStatus(); err != nil {
		c.Error(err)
		internalServerError(c)
		return
	} else {
		renderPowerStatus(c, status)
	}
}

func renderTaskStarted(c *gin.Context, t *power.Task) {
	var res struct {
		ID string `json:"id"`
	}
	res.ID = t.ID()
	c.Header("Location", "/tasks/"+res.ID)
	c.JSON(http.StatusAccepted, &res)
}

func renderPowerStatus(c *gin.Context, status power.Status) {
	var res struct {
		State string `json:"state"`
	}
	res.State = string(status)
	c.JSON(http.StatusOK, &res)
}
