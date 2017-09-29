#!/bin/bash

umask 027

function LOG() {
    content=$1
    echo "[`date '+%Y-%m-%d %H:%M:%S'`] $content"
}

function main() {
    LOG "docker run -it -v $(pwd):/go/src/github.com/cyg2009/MyTestCode -w="/go/src/github.com/cyg2009/MyTestCode" golang:latest go build --ldflags '-extldflags \"-static\"' -a -o processrouter"
    export CGO_ENABLED=0
    export GO_EXTLINK_ENABLED=0
    docker run -it -v $(pwd):/go/src/github.com/cyg2009/MyTestCode -w="/go/src/github.com/cyg2009/MyTestCode" golang:latest go build --ldflags '-extldflags "-static"' -a -o processrouter
    
}

main $*
