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

mkdir -p ${REPO_BASE}
ln -nsf ${PWD} ${REPO_BASE}
for elem in $(find ${PROTOS_LOC} -type f -name '*.proto'); do
    echo -e "Generating .pb.go for ${elem}"
    protoc -I . -I github.com --go_out=${PWD} --go_opt=module=${REPO} ${elem} || true
    protoc -I . -I github.com --go-grpc_out=${PWD} --go-grpc_opt=module=${REPO} ${elem} || true
done
rm ${REPO} && rmdir -p ${REPO_BASE}
# $(go env GOPATH)/bin/goimports -w . || true

