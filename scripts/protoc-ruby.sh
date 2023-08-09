#!/usr/bin/env bash
set -e

PROTOS_LOC=api/v1

pbs=$(find ${PROTOS_LOC} -type f -name '*.proto' -print0 | xargs -0 | tr " " ",")

mkdir -p 'client/ruby/'

echo -n "Generating code from protocol buffers (${pbs})..."
grpc_tools_ruby_protoc -I . -I "$(go env GOMODCACHE)" --ruby_out=. --grpc_out=. ${PROTOS_LOC}/*.proto
mv ${PROTOS_LOC}/*.rb client/ruby/
echo "done"
