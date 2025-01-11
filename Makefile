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

# BEGIN: lint-install .
# http://github.com/tinkerbell/lint-install

GOLINT_VERSION ?= v1.63.4
HADOLINT_VERSION ?= v2.12.0
SHELLCHECK_VERSION ?= v0.10.0
LINT_OS := $(shell uname)
LINT_ARCH := $(shell uname -m)

# shellcheck and hadolint lack arm64 native binaries: rely on x86-64 emulation
ifeq ($(LINT_OS),Darwin)
	ifeq ($(LINT_ARCH),arm64)
		LINT_ARCH=x86_64
	endif
endif

LINT_LOWER_OS  = $(shell echo $(LINT_OS) | tr '[:upper:]' '[:lower:]')
GOLINT_CONFIG:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))/.golangci.yml

lint: out/linters/shellcheck-$(SHELLCHECK_VERSION)-$(LINT_ARCH)/shellcheck out/linters/hadolint-$(HADOLINT_VERSION)-$(LINT_ARCH) out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH)
	out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH) run
	out/linters/hadolint-$(HADOLINT_VERSION)-$(LINT_ARCH) -t info $(shell find . -name "*Dockerfile")
	out/linters/shellcheck-$(SHELLCHECK_VERSION)-$(LINT_ARCH)/shellcheck $(shell find . -name "*.sh")

out/linters/shellcheck-$(SHELLCHECK_VERSION)-$(LINT_ARCH)/shellcheck:
	mkdir -p out/linters
	curl -sSfL https://github.com/koalaman/shellcheck/releases/download/$(SHELLCHECK_VERSION)/shellcheck-$(SHELLCHECK_VERSION).$(LINT_LOWER_OS).$(LINT_ARCH).tar.xz | tar -C out/linters -xJf -
	mv out/linters/shellcheck-$(SHELLCHECK_VERSION) out/linters/shellcheck-$(SHELLCHECK_VERSION)-$(LINT_ARCH)

out/linters/hadolint-$(HADOLINT_VERSION)-$(LINT_ARCH):
	mkdir -p out/linters
	curl -sfL https://github.com/hadolint/hadolint/releases/download/$(HADOLINT_VERSION)/hadolint-$(LINT_OS)-$(LINT_ARCH) > out/linters/hadolint-$(HADOLINT_VERSION)-$(LINT_ARCH)
	chmod u+x out/linters/hadolint-$(HADOLINT_VERSION)-$(LINT_ARCH)

out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH):
	mkdir -p out/linters
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b out/linters $(GOLINT_VERSION)
	mv out/linters/golangci-lint out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH)

# END: lint-install .

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
