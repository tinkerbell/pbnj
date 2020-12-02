// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// ErrBadAccessID indicate an invalid access id
	ErrBadAccessID = errors.New("invalid access id")
	// ErrBadAuthorization indicates an invalid authorization header
	ErrBadAuthorization = errors.New("invalid authorization header")
	// ErrBadDate indicates an invalid date header
	ErrBadDate = errors.New("invalid date header")
	// ErrBadSignature indicates an invalid request signature
	ErrBadSignature = errors.New("invalid request signature")
	// ErrNotSigned indicates an unsigned request
	ErrNotSigned = errors.New("request not signed")
	// ErrTooOld indicates an expired signature
	ErrTooOld = errors.New("request too old")
)

// jsonErrors shoves any error messages into a JSON body, unless the body was already written.
func jsonErrors(c *gin.Context) {
	c.Next()

	var res struct {
		Errors []string `json:"errors"`
	}
	tx := elog.TxFromContext(c)
	for _, e := range c.Errors {
		tx.With("error", e).Error("logging an error")
		res.Errors = append(res.Errors, fmt.Sprintf("%+v", e))
	}

	if c.Writer.Written() || len(c.Errors) == 0 {
		return
	}

	c.JSON(-1, &res)
}

func badRequest(c *gin.Context, errs ...error) {
	for _, err := range errs {
		_ = c.Error(err)
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func internalServerError(c *gin.Context, errs ...error) {
	for _, err := range errs {
		_ = c.Error(err)
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}

func notFound(c *gin.Context, errs ...error) {
	for _, err := range errs {
		_ = c.Error(err)
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func unauthorized(c *gin.Context, errs ...error) {
	for _, err := range errs {
		_ = c.Error(err)
	}
	c.Header("WWW-Authenticate", "APIAuth")
	c.AbortWithStatus(http.StatusUnauthorized)
}
