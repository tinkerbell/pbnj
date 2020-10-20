MAKEFLAGS += --no-builtin-rules
.PHONY: ${binary} run test
.SUFFIXES:

binary := pbnj
all: ${binary}

${binary}: $(shell git ls-files | grep -v -e vendor -e '_test.go' | grep '.go$$' )
	CGO_ENABLED=0 go build -v -ldflags="-X main.GitRev=$(shell git rev-parse --short HEAD)"

ifeq ($(origin GOBIN), undefined)
GOBIN := ${PWD}/bin
export GOBIN
endif

ifeq ($(CI),drone)
run: ${binary}
	${binary}
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ${TEST_ARGS} ./...
else
run: ${binary}
	docker-compose up -d --build pbnj
test:
	docker-compose up -d --build pbnj
endif
