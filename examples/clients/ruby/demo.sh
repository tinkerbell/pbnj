#!/usr/bin/env bash


function setup() {
    gem install grpc grpc-tools > /dev/null 2>&1
    apt update > /dev/null 2>&1; apt -y install jq > /dev/null 2>&1
    for proto in $(ls ../../../api/v1/*.proto); do
        grpc_tools_ruby_protoc -I ../../.. --ruby_out=./lib --grpc_out=./lib $proto 
    done
}

function progress() {
    while true; do
        echo -n .
        sleep 1
        if [ -f "/tmp/setup_done" ]; then
            break
        fi
    done
}

function start_progress() {
    rm -rf /tmp/setup_done
    progress &
}

function end_progress() {
    touch /tmp/setup_done
}

echo -n "setting up environment."
start_progress
setup
end_progress
echo -ne "done\n"

echo -e "calling PBnJ endpoint"
ruby main.rb $1 $2 $3 | jq
