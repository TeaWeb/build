#!/usr/bin/env bash

# run service

export GOPATH=`pwd`/../../

#go build -o teaweb ${GOPATH}/src/github.com/TeaWeb/code/main/main.go
#./teaweb stop
go run ${GOPATH}/src/github.com/TeaWeb/code/main/main.go