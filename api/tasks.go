// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/interfaces/power"
)

func taskStatus(c *gin.Context) {
	task := power.FindTask(c.Param("id"))
	if task == nil {
		notFound(c)
		return
	}

	var timeout time.Duration
	if param := c.Query("timeout"); param != "" {
		d, err := time.ParseDuration(param)
		if err != nil {
			badRequest(c, errors.Wrap(err, "timeout parsing"))
			return
		}
		timeout = d
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-task.Done():
	case <-timer.C:
	}
	renderTaskStatus(c, task)
}

func renderTaskStatus(c *gin.Context, t *power.Task) {
	var res struct {
		Done  bool   `json:"done"`
		Error string `json:"error,omitempty"`
	}
	select {
	case <-t.Done():
		res.Done = true
		res.Error = fmt.Sprintf("%+v", t.Err())
	default:
	}
	c.JSON(http.StatusOK, &res)
}
