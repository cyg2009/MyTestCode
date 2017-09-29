#!/bin/bash

umask 027

function LOG() {
    content=$1
    echo "[`date '+%Y-%m-%d %H:%M:%S'`] $content"
}

function main() {
    LOG "go build --ldflags '-extldflags \"-static\"' -a -o processrouter ./cmd"
    export CGO_ENABLED=0
    export GO_EXTLINK_ENABLED=0
    go build --ldflags '-extldflags "-static"' -a 
}

main $*
