# Functional Tests

This directory holds the functional tests for PBnJ.

## Usage

First populate a `resources.yaml` with BMC resource information.
This example is also found in `resources.example.yaml`.

```yaml
---
server: # if this block is NOT specified than a goroutine is spun up to run a server (see main.go)
  url: 10.4.5.6 
  port: 50051
resources:
- ip: 10.10.10.10 # BMC IP
  username: admin # BMC username
  password: admin # BMC password
  vendor: HP # BMC vendor
  useCases: # use cases for each RPC test
    power: # should correspond to an RPC name
    - happyTests # should correspond to the name of a structure that holds test inputs and expected outputs
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
# under the hood it calls: go run test/main.go -config test/resources.yaml
```
