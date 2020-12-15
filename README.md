# PBNJ
[![Build Status](https://cloud.drone.io/api/badges/tinkerbell/pbnj/status.svg)](https://cloud.drone.io/tinkerbell/pbnj)
![](https://img.shields.io/badge/Stability-Experimental-red.svg)

This service handles BMC interactions.

This repository is [Experimental](https://github.com/packethost/standards/blob/master/experimental-statement.md) meaning that it's based on untested ideas or techniques and not yet established or finalized or involves a radically new and innovative style!
This means that support is best effort (at best!) and we strongly encourage you to NOT use this in production.

## Paths

```
[GIN-debug] GET    /devices/:ip/power        --> github.com/tinkerbell/pbnj/api.powerStatus (5 handlers)
[GIN-debug] POST   /devices/:ip/power        --> github.com/tinkerbell/pbnj/api.powerAction (5 handlers)
[GIN-debug] PATCH  /devices/:ip/boot         --> github.com/tinkerbell/pbnj/api.updateBootOptions (5 handlers)
[GIN-debug] POST   /devices/:ip/bmc          --> github.com/tinkerbell/pbnj/api.bmcAction (5 handlers)
[GIN-debug] PATCH  /devices/:ip/ipmi-lan     --> github.com/tinkerbell/pbnj/api.updateLANConfig (5 handlers)
[GIN-debug] GET    /tasks/:id                --> github.com/tinkerbell/pbnj/api.taskStatus (5 handlers)
[GIN-debug] GET    /                         --> github.com/tinkerbell/pbnj/api.ping (5 handlers)
[GIN-debug] GET    /healthcheck              --> github.com/tinkerbell/pbnj/api.healthcheck (4 handlers)
[GIN-debug] Listening and serving HTTP on :9090
```

## Build


Docker Build
```
docker build -f Dockerfile.dev  .
docker run -it -p 127.0.0.1:9090:9090 <container id>
```

Local
```
# use go get based on what we import
go get ./...
go build
./pbnj
```

Visit http://localhost:9090/healthcheck

## Website

For complete documentation, please visit the Tinkerbell project hosted at [tinkerbell.org](https://tinkerbell.org).
