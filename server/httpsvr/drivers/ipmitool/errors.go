// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"bytes"

	"github.com/pkg/errors"
)

var (
	// ErrConnect mimics the ipmitool response when a connection can not be established.
	ErrConnect = defineError("Unable to establish IPMI v2 / RMCP+ session")
	// ErrStatus mimics the ipmitool response when it is unable to get chasis power status.
	ErrStatus = defineError("Unable to get Chassis Power Status")
	// ErrChannelCipherSuites mimics the ipmitool response when it there is a cypher suites mismatch.
	ErrChannelCipherSuites = defineError("Unable to Get Channel Cipher Suites")
)

var knownErrors = make(map[string]error)

func defineError(msg string) error {
	err, ok := knownErrors[msg]
	if ok {
		logger.Panicf("known error %q already defined!", msg)
	} else {
		err = errors.New(msg)
		knownErrors[msg] = err
	}
	return err
}

func lookupError(msg []byte) error {
	if err, ok := knownErrors[string(msg)]; ok {
		return err
	}
	str := string(msg)
	return errors.New(str)
}

func parseStderrLine(line []byte) error {
	if len(line) == 0 {
		return nil
	}
	line = bytes.TrimPrefix(line, []byte("> "))
	line = bytes.TrimPrefix(line, []byte("Error: "))
	return lookupError(line)
}
