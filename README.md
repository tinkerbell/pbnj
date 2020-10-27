# PBNJ

[![Build Status](https://cloud.drone.io/api/badges/tinkerbell/pbnj/status.svg)](https://cloud.drone.io/tinkerbell/pbnj)
![](https://img.shields.io/badge/Stability-Experimental-red.svg)

This service handles BMC interactions.

This repository is [Experimental](https://github.com/packethost/standards/blob/master/experimental-statement.md) meaning that it's based on untested ideas or techniques and not yet established or finalized or involves a radically new and innovative style!
This means that support is best effort (at best!) and we strongly encourage you to NOT use this in production.

## Usage

### Container

Build
```
make image
```

Run
```
# default port is 50051
make run-image
```

### Local

Build
```
# builds the binary and puts it in ./bin/
make build
```

Run
```
# default port is 50051; does a `go run` of the code base
make run-server
```
