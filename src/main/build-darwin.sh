#!/usr/bin/env bash

# build project for darwin (Mac OS X)

. utils.sh

export GOPATH=`pwd`/../../
export GO111MODULE=off
export GOOS=darwin
export GOARCH=amd64

build
