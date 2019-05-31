#!/bin/bash

set -eux

export GOPATH="$(pwd)/.gobuild"
SRCDIR="${GOPATH}/src/github.com/trendyol/easy-rollback"
export GO111MODULE=on
[ -d ${GOPATH} ] && rm -rf ${GOPATH}
mkdir -p ${GOPATH}/{src,pkg,bin}
mkdir -p ${SRCDIR}
cp -r ./* ${SRCDIR}
(
    echo ${GOPATH}
    cd ${SRCDIR}
    go get . 
    go build .
    go install .
)
