#!/usr/bin/env bash

# build project for darwin (Mac OS X)

	CWD=$(dirname $0)

. $CWD/utils.sh

export GOPATH=""
export GO111MODULE=on
export GOOS=darwin
export GOARCH=amd64

build
