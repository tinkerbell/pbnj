// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"bytes"
)

func eachLine(buf *bytes.Buffer, p []byte, fn func([]byte)) {
	i := bytes.IndexByte(p, '\n')
	if buffered := buf.Len(); buffered != 0 {
		buf.Write(p)

		if i == -1 {
			return
		}
		i = buffered + i
		p = buf.Next(i + 1)
	}
	for i != -1 {
		fn(p[:i])
		p = p[i+1:]
		i = bytes.IndexByte(p, '\n')
	}
	if len(p) > 0 {
		buf.Write(p)
	}
}

type funcWriter func([]byte) (int, error)

func (f funcWriter) Write(p []byte) (int, error) { return f(p) }
