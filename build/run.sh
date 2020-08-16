#!/usr/bin/env bash

# run service, you can pass options to this shell, such as './run.sh pprof'

CWD=$(dirname $0)

export GOPATH=""
export GO111MODULE=on

go run ${CWD}/../cmd/teaweb/main.go $1 $2 $3 $4