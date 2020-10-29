# PBNJ
[![Build Status](https://cloud.drone.io/api/badges/tinkerbell/pbnj/status.svg)](https://cloud.drone.io/tinkerbell/pbnj)
![](https://img.shields.io/badge/Stability-Experimental-red.svg)

This service handles BMC interactions.

This repository is [Experimental](https://github.com/packethost/standards/blob/master/experimental-statement.md) meaning that it's based on untested ideas or techniques and not yet established or finalized or involves a radically new and innovative style!
This means that support is best effort (at best!) and we strongly encourage you to NOT use this in production.

## Paths

```
GET    /devices/:ip/power        --> api.powerStatus (5 handlers)
POST   /devices/:ip/power        --> api.powerAction (5 handlers)
PATCH  /devices/:ip/boot         --> api.updateBootOptions (5 handlers)
POST   /devices/:ip/bmc          --> api.bmcAction (5 handlers)
PATCH  /devices/:ip/ipmi-lan     --> api.updateLANConfig (5 handlers)
GET    /tasks/:id                --> api.taskStatus (5 handlers)
GET    /                         --> api.ping (5 handlers)
GET    /healthcheck              --> api.healthcheck (4 handlers)
GET    /redfish/*                --> api.redfish.Proxy (1 handler)
POST   /redfish/*                --> api.redfish.Proxy (1 handler)
Listening and serving HTTP on :9090
```


## Build & Run

Docker Build
```
docker build -f Dockerfile.dev  .
docker run -it -p 127.0.0.1:9090:9090 <container id>
```

## Examples

Setup:
``bash
// setup the environment variables for the BMC you wish to control:
IP=ip-address
IPMI_USER=username
IPMI_PASS=password
```

Power Status:
```bash
curl \
  -H "X-IMPI-Username: ${IPMI_USER}" \
  -H "X-IMPI-Password: ${IPMI_PASS}" \
  "http://localhost:9090/devices/${IP}/power"

```

Passthru Command:
```bash
curl \
  -H "X-IMPI-Username: ${IPMI_USER}" \
  -H "X-IMPI-Password: ${IPMI_PASS}" \
  -d '{
  "action": "command",
  "command": "lan set 1 ipaddr 192.168.1.1"
}' \
  "http://localhost:9090/devices/${IP}/bmc"
```


Local
```
# use go get based on what we import
go get ./...
go build
./pbnj
```

Visit http://localhost:9090/healthcheck
