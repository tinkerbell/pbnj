# PBNJ

![For each commit and PR](https://github.com/tinkerbell/pbnj/workflows/For%20each%20commit%20and%20PR/badge.svg)
![stability](https://img.shields.io/badge/Stability-Experimental-red.svg)

> This repository is [Experimental](https://github.com/packethost/standards/blob/main/experimental-statement.md) meaning that it's based on untested ideas or techniques and not yet established or finalized or involves a radically new and innovative style!
> This means that support is best effort (at best!) and we strongly encourage you to NOT use this in production.

## Description

This service handles BMC interactions.

- machine and BMC power on/off/reset
- setting next boot device
- user management
- setting BMC network source

The gRPC PBnJ server listens by default on port 50051.
This can be started with `pbnj server`.
Use `pbnj server --help` for more runtime details.

## Usage

### Container

Build

```bash
make image
```

Run

```bash
# default gRPC port is 50051
make run-image
```

### Local

Build

```bash
# builds the binary and puts it in ./bin/
make build
```

Run

```bash
# default gRPC port is 50051; does a `go run` of the code base
make run-server
```

## Authorization

Documentation on enabling authorization can be found [here](docs/Authorization.md).

## Contributing

See the contributors guide [here](CONTRIBUTING.md).

## Website

For complete documentation, please visit the Tinkerbell project hosted at [tinkerbell.org](https://tinkerbell.org).
