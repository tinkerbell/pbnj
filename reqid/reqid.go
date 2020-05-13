// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package reqid

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/log"
)

// id is the type of value stored in the Contexts.
type id string

// key is the key for id values in Contexts.
// It is unexported; clients use reqid.NewContext and reqid.FromContext instead of using this key directly.
const key = id("reqid")

var (
	logger log.Logger
)

func SetupLogging(l log.Logger) {
	logger = l.Package("reqid")
}

// WithID returns a new Context that carries value id
func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key, id)
}

// Set is used to store the id exclusively for this gin context.
func Set(ctx *gin.Context, id string) {
	ctx.Set(string(key), id)
}

// FromContext returns the id value stored in ctx, if any.
func FromContext(ctx context.Context) string {
	if ctx == nil {
		panic("nil context")
	}
	if v := ctx.Value(key); v != nil {
		return v.(string)
	}
	if v := ctx.Value(string(key)); v != nil {
		return v.(string)
	}
	if c, ok := ctx.(*gin.Context); ok {
		id := New()
		c.Set(string(key), id)
		return id
	}

	return ""
}

// New creates a new id
func New() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		err = errors.Wrap(err, "failed to read random bytes")
		logger.Error(err)
		return ""
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
}
