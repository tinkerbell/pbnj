# Code Structure

The goal of this document is to help you understand and navigate the PBnJ code base.

```bash
$ tree -d -L 3
.
├── api
│   └── v1
├── client
├── cmd
├── docs
├── examples
│   └── clients
│       └── ruby
├── grpc
│   ├── oob
│   │   ├── bmc
│   │   └── machine
│   ├── persistence
│   ├── rpc
│   └── taskrunner
├── pkg
│   ├── healthcheck
│   ├── http
│   ├── logging
│   ├── metrics
│   ├── oob
│   ├── repository
│   ├── task
│   └── zaplog
├── scripts
└── test
    └── runner
```

## Table of Contents

- [Code Structure](#code-structure)
  - [Table of Contents](#table-of-contents)
    - [_api/_](#api)
    - [_cmd/_](#cmd)
    - [_pkg/_](#pkg)
    - [_scripts/_](#scripts)
    - [_grpc/_](#grpc)
    - [_test/_](#test)

### _api/_

The _api_ directory holds versioned directories of protocol buffer files (.proto) and generated code.
The code generated from the protocol buffers can be produced in 2 ways.

1. Using local dependencies

   ```bash
   # install the dependencies locally
   make pbs-install-deps

   # generate the protobuf code
   make pbs
   ```

2. Using Docker

   ```bash
   make pbs-docker-image
   ```

### _cmd/_

The _cmd_ directory is the entrypoint code for the PBnJ server and client.

### _pkg/_

The _pkg_ directory contains generic utility type Go packages for PBnJ.

### _scripts/_

The _scripts_ directory contains helper scripts for things like protocol buffer code generation, container image running, and running the server locally.

### _grpc/_

The _grpc_ directory contains the gRPC server implementations of PBnJ.

### _test/_

The _test_ directory contains code for running function/integration tests against live hardware.
Automating this type of testing is difficult so running this code requires special attention.
See the [README](../test/README.md) for details.
