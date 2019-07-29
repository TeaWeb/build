#!/usr/bin/env bash

# run service

export GOPATH=`pwd`/../../
export GO111MODULE=off

go run -ldflags="-s -w" ${GOPATH}/src/github.com/TeaWeb/code/main/main.go