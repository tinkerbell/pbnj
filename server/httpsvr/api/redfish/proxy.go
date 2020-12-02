// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package redfish

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

var (
	proxyTLS    = &httputil.ReverseProxy{Director: director}
	proxyNonTLS = &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	redfishCTX struct{}
)

func director(req *http.Request) {
	ctx := req.Context().Value(redfishCTX).(*gin.Context)

	req.URL.Scheme = "https"
	if req.Header.Get("X-REDFISH-SCHEME") == "http" {
		req.URL.Scheme = "http"
	}

	username := req.Header.Get("X-IPMI-Username")
	password := req.Header.Get("X-IPMI-Password")
	req.SetBasicAuth(username, password)

	req.Host = ctx.Param("ip")
	req.RequestURI = "/redfish" + ctx.Param("redfish")
	req.URL.Host = ctx.Param("ip")
	req.URL.Path = "/redfish" + ctx.Param("redfish")
}

// redfishProxy is a handler for the ANY /redfish/* endpoint
// Proxies /device/{ip}/redfish/{path} to {ip}/redfish/{path}
func Proxy(c *gin.Context) {
	proxy := proxyTLS
	if c.Request.Header.Get("X-REDFISH-TLS-VERIFY") == "false" {
		proxy = proxyNonTLS
	}
	reqctx := context.WithValue(c, redfishCTX, c)
	proxy.ServeHTTP(c.Writer, c.Request.WithContext(reqctx))
}
