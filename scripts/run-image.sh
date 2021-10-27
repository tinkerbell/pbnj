#!/usr/bin/env bash

# why this script you may ask?
# need to be able to clean up the docker image from a crtl-c
# and it's not straightforward in a Makefile.

# trap ctrl-c and call ctrl_c()
trap ctrl_c INT

function ctrl_c() {
    docker rm -f "${CONTAINER_ID}"
    exit 0
}

CONTAINER_ID=$(docker run -d -e PNBJ_LOGLEVEL=debug -p "${PBNJ_PORT:-50051}":50051 -p 8080:8080 pbnj:local)
docker logs -f "${CONTAINER_ID}" 2>&1 | jq -R 'fromjson? | select(type == "object")'