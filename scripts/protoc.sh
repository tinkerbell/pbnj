#!/usr/bin/env bash
#
# protoc.sh uses the local protoc if installed, otherwise
# docker will be used with a complete environment provided
# by https://github.com/jaegertracing/docker-protobuf.
# Alternative images like grpc/go are very dated and do not
# include the needed plugins and includes.
#
set -e

REPO=github.com/tinkerbell/pbnj
REPO_BASE=$(dirname ${REPO})
PROTOS_LOC=api/v1
OUT_LOC=api/v1

function installDeps {
    if ! which protoc 2>&1 >/dev/null; then
        echo 'Installing protoc...' >&2
        if which apt 2>&1 >/dev/null; then
            apt update; apt install -y protobuf-compiler
        elif which yum 2>&1 >/dev/null; then
            yum -y install protobuf-compiler
        else
            echo 'Unknown package manager' >&2
            exit 1
        fi
        echo 'Protoc installed' >&2
    else
        echo 'Protoc already installed!' >&2
    fi
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.1
    go get google.golang.org/grpc@v1.32.0
    go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
    go get golang.org/x/tools/cmd/goimports
    go mod tidy
}

if [[ "$1" == "deps" ]]; then
    installDeps
    exit 0
fi

mkdir -p ${REPO_BASE}
ln -nsf ${PWD} ${REPO_BASE}
for elem in $(find ${PROTOS_LOC} -type f -name '*.proto'); do
    echo -e "Generating .pb.go for ${elem}"
    protoc -I . -I github.com --go_out=${PWD} --go_opt=module=${REPO} ${elem} || true
    protoc -I . -I github.com --go-grpc_out=${PWD} --go-grpc_opt=module=${REPO} ${elem} || true
done
rm -f ${REPO}
rmdir -p ${REPO_BASE}
$(go env GOPATH)/bin/goimports -w . || true
