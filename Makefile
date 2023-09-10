REPO:=github.com/tinkerbell/pbnj
REPO_BASE:=$(shell dirname ${REPO})
PROTOS_LOC:=v2/protos
BINARY:=pbnj
OSFLAG:= $(shell go env GOHOSTOS)
GIT_COMMIT:=$(shell git rev-parse --short HEAD)
BUILD_ARGS:=GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -extldflags "-static"'
PROTOBUF_BUILDER_IMG:=pbnj-protobuf-builder

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: darwin
darwin: ## complie for darwin
	GOOS=darwin ${BUILD_ARGS} -o bin/${BINARY}-darwin-amd64 main.go

.PHONY: linux
linux: ## complie for linux
	GOOS=linux ${BUILD_ARGS} -o bin/${BINARY}-linux-amd64 main.go

.PHONY: build
build: ## compile the binary for the native OS
ifeq (${OSFLAG},linux)
	@$(MAKE) linux
else
	@$(MAKE) darwin
endif

.PHONY: image
image: ## make the Container Image
	docker build -t pbnj:local .

##@ Development

.PHONY: test
test: ## run tests
	go test -v -covermode=count ./...

.PHONY: test-ci
test-ci: ## run tests for ci and codecov
	go test -coverprofile=coverage.txt ./...

.PHONY: test-functional
test-functional: ## run functional tests
	go run test/main.go -config test/resources.yaml

.PHONY: goimports-ci
goimports-ci: ## run goimports for ci
	go get golang.org/x/tools/cmd/goimports
	test -z "$(shell ${GOBIN}/goimports -d -l ./| tee /dev/stderr)"

.PHONY: goimports
goimports: ## run goimports
	@echo be sure goimports is installed
	goimports -w ./

.PHONY: cover
cover: ## Run unit tests with coverage report
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out
	rm -rf cover.out

.PHONY: buf-lint
buf-lint:  ## run linting
	@echo be sure buf is installed: https://buf.build/docs/installation
	buf check lint

PHONY: run-server
run-server: ## run server locally
ifeq (, $(shell which jq))
	go run ./cmd/${BINARY} server
else
	scripts/run-server.sh
endif

.PHONY: pbs
pbs: ## locally generate go stubs from protocol buffers
	scripts/protoc.sh

.PHONY: pbs-install-deps
pbs-install-deps: ## locally install dependencies in order to generate go stubs from protocol buffers
	scripts/protoc.sh deps

.PHONY: pbs-docker
pbs-docker: pbs-docker-image ## generate go stubs from protocol buffers in a container
	docker run -it --rm -v ${PWD}:/code -w /code ${PROTOBUF_BUILDER_IMG} scripts/protoc.sh

.PHONY: pbs-docker-ruby
pbs-docker-ruby:  ## generate ruby stubs from protocol buffers in a container
	docker build -t ${PROTOBUF_BUILDER_IMG}.ruby -f scripts/Dockerfile.pbbuilder.ruby .
	docker run -it --rm -v ${PWD}:/code -w /code ${PROTOBUF_BUILDER_IMG}.ruby scripts/protoc-ruby.sh

.PHONY: pbs-docker-image
pbs-docker-image: ## generate container image for building protocol buffers
	docker build -t ${PROTOBUF_BUILDER_IMG} -f scripts/Dockerfile.pbbuilder .

.PHONY: run-image
run-image: ## run PBnJ container image
	scripts/run-image.sh

# BEGIN: lint-install github.com/tinkerbell/pbnj
# http://github.com/tinkerbell/lint-install

.PHONY: lint
lint: _lint

LINT_ARCH := $(shell uname -m)
LINT_OS := $(shell uname)
LINT_OS_LOWER := $(shell echo $(LINT_OS) | tr '[:upper:]' '[:lower:]')
LINT_ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# shellcheck and hadolint lack arm64 native binaries: rely on x86-64 emulation
ifeq ($(LINT_OS),Darwin)
	ifeq ($(LINT_ARCH),arm64)
		LINT_ARCH=x86_64
	endif
endif

LINTERS :=
FIXERS :=

GOLANGCI_LINT_CONFIG := $(LINT_ROOT)/.golangci.yml
GOLANGCI_LINT_VERSION ?= v1.53.3
GOLANGCI_LINT_BIN := $(LINT_ROOT)/out/linters/golangci-lint-$(GOLANGCI_LINT_VERSION)-$(LINT_ARCH)
$(GOLANGCI_LINT_BIN):
	mkdir -p $(LINT_ROOT)/out/linters
	rm -rf $(LINT_ROOT)/out/linters/golangci-lint-*
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINT_ROOT)/out/linters $(GOLANGCI_LINT_VERSION)
	mv $(LINT_ROOT)/out/linters/golangci-lint $@

LINTERS += golangci-lint-lint
golangci-lint-lint: $(GOLANGCI_LINT_BIN)
	find . -name go.mod -execdir "$(GOLANGCI_LINT_BIN)" run -c "$(GOLANGCI_LINT_CONFIG)" \;

FIXERS += golangci-lint-fix
golangci-lint-fix: $(GOLANGCI_LINT_BIN)
	find . -name go.mod -execdir "$(GOLANGCI_LINT_BIN)" run -c "$(GOLANGCI_LINT_CONFIG)" --fix \;

.PHONY: _lint $(LINTERS)
_lint: $(LINTERS)

.PHONY: fix $(FIXERS)
fix: $(FIXERS)

# END: lint-install github.com/tinkerbell/pbnj


##@ Clients

.PHONY: ruby-client-demo
ruby-client-demo: image ## run ruby client demo
	# make ruby-client-demo host=10.10.10.10 user=ADMIN pass=ADMIN
	docker run -d --name pbnj pbnj:local
	docker run -it --rm --net container:pbnj -v ${PWD}:/code -w /code/examples/clients/ruby --entrypoint /bin/bash ruby /code/examples/clients/ruby/demo.sh ${host} ${user} ${pass}
	docker rm -f pbnj

.PHONY: evans
evans: ## run evans grpc client
	evans --path $$(go env GOMODCACHE) --path . --proto $$(find api/v1 -type f -name '*.proto'| xargs | tr " " ",") -p "$${PBNJ_PORT:-50051}" repl
