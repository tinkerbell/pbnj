// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	maxRequestAge = 15 * time.Minute
)

func authorize(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		if gin.Mode() == gin.DebugMode {
			logger.Warning("ignoring unsigned request because we're in debug mode")
			return
		}

		unauthorized(c, ErrNotSigned)
		return
	}

	// Support for the api-auth gem
	if auth := strings.TrimPrefix(authHeader, "APIAuth "); len(auth) < len(authHeader) {
		if !validateTimestamp(c, maxRequestAge) {
			return
		}

		key, sig, err := splitAPIAuth(auth)
		if err != nil {
			unauthorized(c, ErrBadSignature)
			return
		}

		validateSignature(c, sig, key)
		return
	}

	unauthorized(c, ErrNotSigned)
}

var (
	primaryAccessID     = os.Getenv("ACCESS_ID")
	primaryAccessSecret = os.Getenv("ACCESS_SECRET")
)

func lookupKey(accessID string) []byte {
	if accessID == primaryAccessID {
		return []byte(primaryAccessSecret)
	}
	return nil
}

func splitAPIAuth(auth string) (key, sig []byte, err error) {
	i := strings.IndexByte(auth, ':')
	if i == -1 {
		return nil, nil, ErrBadAuthorization
	}

	key = lookupKey(auth[:i])
	if key == nil {
		return nil, nil, ErrBadAccessID
	}

	sig, err = base64.StdEncoding.DecodeString(auth[i+1:])
	if err != nil {
		return nil, nil, errors.WithMessage(ErrBadSignature, "decoding api auth string")
	}
	return key, sig, nil
}

func validateSignature(c *gin.Context, sig, key []byte) bool {
	var (
		method      = c.Request.Method
		uri         = c.Request.RequestURI
		contentType = c.Request.Header.Get("Content-Type")
		contentMD5  = c.Request.Header.Get("Content-MD5")
		timestamp   = c.Request.Header.Get("Date")
	)

	mac := hmac.New(sha1.New, key)
	fmt.Fprintf(mac, "%s,%s,%s,%s,%s", method, contentType, contentMD5, uri, timestamp)

	if expected := mac.Sum(nil); !hmac.Equal(sig, expected) {
		unauthorized(c, ErrBadSignature)
		return false
	}
	return true
}

var allowedDateFormats = []string{
	http.TimeFormat,
	time.RFC1123Z,
	time.RFC1123,
	time.RFC850,
	time.ANSIC,
}

func validateTimestamp(c *gin.Context, maxAge time.Duration) bool {
	var (
		timestamp = c.Request.Header.Get("Date")

		t   time.Time
		err error
	)
	for _, format := range allowedDateFormats {
		t, err = time.Parse(format, timestamp)
		if err == nil {
			break
		}
	}
	if err != nil {
		badRequest(c, ErrBadDate)
		return false
	}

	if age := time.Since(t); age > maxAge {
		unauthorized(c, ErrTooOld)
		return false
	}
	return true
}
