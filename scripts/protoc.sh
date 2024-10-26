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
PROTOS_LOC=api/v1
PROTOC_VERSION=3.13.0

function installDeps {
    if ! which protoc &>/dev/null; then
        echo 'Installing protoc...' >&2
        if ! which unzip &>/dev/null; then
            if which apt &>/dev/null; then
                apt update; apt install -y zip
            elif which yum &>/dev/null; then
                yum -y install zip
            else
                echo 'Unknown package manager' >&2
                exit 1
            fi
        fi
        PB_REL="https://github.com/protocolbuffers/protobuf/releases"
        curl -LO $PB_REL/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip
        unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local
        echo 'Protoc installed' >&2
    else
        echo 'Protoc already installed!' >&2
    fi
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.1
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
    go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@v0.3.2
}

if [[ "$1" == "deps" ]]; then
    installDeps
    exit 0
fi
pbs=$(find ${PROTOS_LOC} -type f -name '*.proto' -print0 | xargs -0 | tr " " ",")

echo -n "Generating code from protocol buffers (${pbs})..."
protoc -I . -I "$(go env GOMODCACHE)" --go_out=. --go_opt=module=${REPO} ${PROTOS_LOC}/*.proto
protoc -I . -I "$(go env GOMODCACHE)" --govalidators_out=. --go-grpc_out=. --go-grpc_opt=module=${REPO} ${PROTOS_LOC}/*.proto
mv ${REPO}/${PROTOS_LOC}/*.go ${PROTOS_LOC}/
rm -rf github.com
echo "done"
