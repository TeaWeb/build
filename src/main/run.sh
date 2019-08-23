#!/usr/bin/env bash

# run service

export GOPATH=`pwd`/../../
export GO111MODULE=off

go run ${GOPATH}/src/github.com/TeaWeb/code/main/main.go