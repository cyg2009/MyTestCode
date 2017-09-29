#!/bin/bash

umask 027

function LOG() {
    content=$1
    echo "[`date '+%Y-%m-%d %H:%M:%S'`] $content"
}

function main() {
    LOG "go build --ldflags '-extldflags \"-static\"' -a -o funcRouter"
    CGO_ENABLED=0 GO_EXTLINK_ENABLED=0 go build --ldflags '-extldflags "-static"' -a -o funcRouter  
    docker build -t my:1.0 .
}

main $*
