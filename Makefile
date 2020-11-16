REPO:=github.com/tinkerbell/pbnj
REPO_BASE:=$(shell dirname ${REPO})
PROTOS_LOC:=v2/protos
BINARY:=pbnj
OSFLAG:= $(shell go env GOHOSTOS)
GIT_COMMIT:=$(shell git rev-parse --short HEAD)
BUILD_ARGS:=GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -extldflags "-static"'
PROTOBUF_BUILDER_IMG:=pbnj-protobuf-builder

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## run tests
	go test -v -covermode=count ./...

.PHONY: test-ci
test-ci: ## run tests for ci and codecov
	go test -coverprofile=coverage.txt ./...

.PHONY: test-functional
test-functional: ## run functional tests
	go test -v ./test/... --tags=functional -config 'resources.yaml'

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

.PHONY: lint
lint:  ## run linting
	@echo be sure golangci-lint is installed: https://golangci-lint.run/usage/install/
	golangci-lint run

.PHONY: buf-lint
buf-lint:  ## run linting
	@echo be sure buf is installed: https://buf.build/docs/installation
	buf check lint

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

PHONY: run-server
run-server: ## run server locally
ifeq (, $(shell which jq))
	go run main.go server
else
	scripts/run-server.sh
endif

.PHONY: pbs
pbs: ## locally generate go stubs from protocol buffers
	scripts/protoc.sh

.PHONY: pbs-install-deps
pbs: ## locally install dependencies in order to generate go stubs from protocol buffers
	scripts/protoc.sh deps

.PHONY: pbs-docker
pbs-docker: pbs-docker-image ## generate go stubs from protocol buffers in a container
	docker run -it --rm -v ${PWD}:/code -w /code ${PROTOBUF_BUILDER_IMG} scripts/protoc.sh

.PHONY: pbs-docker-image
pbs-docker-image: ## generate container image for building protocol buffers 
	docker build -t ${PROTOBUF_BUILDER_IMG} -f scripts/Dockerfile.pbbuilder .

.PHONY: image
image: ## make the Container Image
	docker build -t pbnj:local . 

.PHONY: run-image
run-image: ## run PBnJ container image
	docker run -it --rm -e PNBJ_LOGLEVEL=debug -p 9090:9090 pbnj:local
