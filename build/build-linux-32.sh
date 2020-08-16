#!/usr/bin/env bash

# build project for linux 32 bit

CWD=$(dirname $0)

. $CWD/utils.sh

export GOPATH=""
export GO111MODULE=on
export GOOS=linux
export GOARCH=386

build
