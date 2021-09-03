# Contributor Guide

Welcome to PBnJ! We are really excited to have you.
Please use the following guide on your contributing journey.
Thanks for contributing!

## Table of Contents

- [Context](#Context)
- [Architecture](#Architecture)
  - [Design Docs](#Design-Docs)
  - [Code Structure](#Code-Structure)
- [Prerequisites](#Prerequisites)
  - [DCO Sign Off](#DCO-Sign-Off)
  - [Code of Conduct](#Code-of-Conduct)
  - [Setting up your development environment](#Setting-up-your-development-environment)
- [Development](#Development)
  - [Building](#Building)
  - [Unit testing](#Unit-testing)
  - [Linting](#Linting)
  - [Functional testing](#Functional-testing)
  - [Running PBnJ locally](#Running-PBnJ-locally)
- [Pull Requests](#Pull-Requests)
  - [Branching strategy](#Branching-strategy)
  - [Quality](#Quality)
    - [CI](#CI)
    - [Code coverage](#Code-coverage)

---

## Context

PBnJ is a service that handles interacting with many different BMC models.
It is part of the [Tinkerbell stack](https://tinkerbell.org) and provides the glue for machine provisioning by enabling machine restarts and setting next boot devices.

## Architecture

### Design Docs

Details and diagrams for PBnJ are found [here](docs/design/DESIGN.md).

### Code Structure

Details on PBnJ's code structure is found [here](docs/CODE_STRUCTURE.md)

## Prerequisites

### DCO Sign Off

Please read and understand the DCO found [here](docs/DCO.md).

### Code of Conduct

Please read and understand the code of conduct found [here](https://github.com/tinkerbell/.github/blob/main/CODE_OF_CONDUCT.md).

### Setting up your development environment

1. Install Go

   PBnJ requires [Go 1.15](https://golang.org/dl/) or later.

2. Install Docker

   PBnJ uses Docker for protocol buffer code generation, container image builds and for the Ruby client example.
   Most versions of Docker will work.

> The items below are nice to haves, but not hard requirements for development

1. Install Evans

   [Evans](https://github.com/ktr0731/evans#installation) is a gRPC client.
   There is a make target (`make evans`) that will setup the connection to the locally running server (`make run-server`)

2. Install golangci-lint

   [golangci-lint](https://golangci-lint.run/usage/install/) is used in CI for lint checking and should be run locally before creating a PR.

3. Install jq

   [jq](https://stedolan.github.io/jq/download/) is used when running the PBnJ server to pretty print the logs.

4. Install buf

   [buf](https://buf.build/docs/installation) is used for linting the protocol buffers.

## Development

### Building

To build PBnJ using your local Go environment, run:

```bash
# will build a binary in the native OS format
make build

# you can specify the OS binary format
make linux
make darwin

# Built binaries can be found in ./bin/
```

To build the PBnJ container image, run:

```bash
make image

# Built image will be named pbnj:local
```

### Unit testing

To execute the unit tests, run:

```bash
make test

# to get code coverage numbers, run:
make cover
```

### Linting

To execute linting, run:

```bash
# runs golangci-lint
make lint

# runs goimports
make goimports

# CI runs go vet, there's no make target for that since its built in.
go vet ./...
```

### Functional testing

The _test/_ directory holds the code for running functional tests against real hardware.
See the [README](test/README.md) in the test directory for details.

### Running PBnJ locally

Locally, PBnJ can be run in two different ways.

1. Local direct

   ```bash
   make run-server
   ```

2. Local with Docker

   ```bash
   # build the container image
   make image

   # run the image
   make run-image
   ```

   To start a PBnJ client, after getting PBnJ running with one of the two options above, run:

   ```bash
   make evans
   ```

## Pull Requests

### Branching strategy

PBnJ uses a fork and pull request model.
See this [doc](https://guides.github.com/activities/forking/) for more details.

### Quality

#### CI

PBnJ uses GitHub Actions for CI.
The workflow is found in [.github/workflows/ci.yaml](.github/workflows/ci.yaml).
It is run for each commit and PR.
The container image building only happens once a PR is merged into the main line.

#### Code coverage

PBnJ does run code coverage with each PR.
Coverage thresholds are not currently enforced.
It is always nice and very welcomed to add tests and keep or increase the code coverage percentage.

### Pre PR Checklist

This checklist is a helper to make sure there's no gotchas that come up when you submit a PR.

- [ ] You've reviewed the [code of conduct](#Code-of-Conduct)
- [ ] All commits are DCO signed off
- [ ] Code is [formatted and linted](#Linting)
- [ ] Code [builds](#Building) successfully
- [ ] All tests are [passing](#Unit-testing)
- [ ] Code coverage [percentage](#Code-coverage). (main line is the base with which to compare)
