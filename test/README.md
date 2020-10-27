# Functional Tests

This directory holds the functional tests for PBnJ.

## Usage

First populate a `resources.yaml` with BMC resource information. This example is also found in `resources.example.yaml`. This file needs to exist in the top level of the `test` directory.
```yaml
---
server: # if this block is NOT specified than a goroutine is spun up to run a server (see common_test.go:TestMain function)
  url: 10.4.5.6 
  port: 50051
resources:
- ip: 10.10.10.10 # BMC IP
  username: admin # BMC username
  password: admin # BMC password
  vendor: HP # BMC vendor
  useCases: # use cases for each RPC test
    power: # should correspond to an RPC name
    - happyTests # should correspond to the name of a structure that holds test inputs and expected outputs; see table driven tests pattern
      sadTests # as many of this as needed can be added
    device:
    - happyTests
      sadTests
- ip: 10.5.5.5
  username: ADMIN
  password: ADMIN
  vendor: supermicrox10
  useCases: 
    power: 
    - happyTests
      sadTests
    device:
    - happyTests
      sadTests
```

Run the tests
```bash
make test-functional
# under the hood it calls: go test -v ./test/... --tags=functional -config 'resources.yaml'
```

## Development

* The functional tests use the standard `Go` testing library.
* Test files should follow normal `go test` patterns and be named with `_test.go`
* Test files should contain a [build tag](https://golang.org/cmd/go/#hdr-Build_constraints) of `// +build functional` at the start of the file.
* Currently, the functional tests are organized with a file per protocol buffer services (`BMC`, `Machine`, and `User`).
* Look at `machine_test.go` -> `runClient` for one option of running the gRPC client.
* Be aware that the `TestMain` function in `common_test.go` is used for all tests in this package. It parses the `resources.yaml` and starts the gRPC server.
* The gRPC server logs will be streamed to stdout. To disable this and only print the server logs for failed runs, remove the `-v` from the `go test` command in the `Makefile`
