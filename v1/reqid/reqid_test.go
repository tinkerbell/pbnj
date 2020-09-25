// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package reqid

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert := require.New(t)
	assert.NotEqual(New(), New())
}

func TestWithID(t *testing.T) {
	assert := require.New(t)
	ctx := context.TODO()
	assert.Nil(ctx.Value(key))
	ctx = WithID(ctx, "foo")
	assert.Equal(ctx.Value(key), "foo")
}

func TestFromContext(t *testing.T) {
	assert := require.New(t)

	ctx := context.TODO()
	assert.Equal(FromContext(ctx), "")
	assert.Equal(FromContext(ctx), FromContext(ctx))
	ctx = context.WithValue(ctx, key, "foo")
	assert.Equal(FromContext(ctx), "foo")
	assert.Equal(FromContext(ctx), FromContext(ctx))

	gctx := &gin.Context{}
	assert.NotEqual(FromContext(gctx), "")
	assert.Equal(FromContext(gctx), FromContext(gctx))
}
