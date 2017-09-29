#!/bin/bash

umask 027

function LOG() {
    content=$1
    echo "[`date '+%Y-%m-%d %H:%M:%S'`] $content"
}

function main() {
    LOG "docker run -v $(pwd):/go/src/github.com/cyg2009/MyTestCode -w="/go/src/github.com/cyg2009/MyTestCode" golang:1.9 go build --ldflags '-extldflags \"-static\"' -a -o funcRouter"
    export CGO_ENABLED=0
    export GO_EXTLINK_ENABLED=0
    docker run -e CGO_ENABLED=0 -e GO_EXTLINK_ENABLED=0 -v $(pwd):/go/src/github.com/cyg2009/MyTestCode -w="/go/src/github.com/cyg2009/MyTestCode" golang:1.9  go build --ldflags '-extldflags "-static"' -a -o funcRouter  
}

main $*
