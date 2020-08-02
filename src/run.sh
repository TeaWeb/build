#!/usr/bin/env bash

# run service, you can pass options to this shell, such as './run.sh pprof'


export GO111MODULE=on

go run main.go $1 $2 $3 $4