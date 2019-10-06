#!/usr/bin/env bash

# run service, you can pass options to this shell, such as './run.sh pprof'

export GOPATH=`pwd`/../../
export GO111MODULE=off

go run ${GOPATH}/src/github.com/TeaWeb/code/main/main.go $1 $2 $3 $4