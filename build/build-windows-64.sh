#!/usr/bin/env bash

# build project for linux 64 bit

CWD=$(dirname $0)

. $CWD/utils.sh

export GOPATH=""
export GO111MODULE=on
export GOOS=windows
export GOARCH=amd64

build
