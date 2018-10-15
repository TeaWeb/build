#!/usr/bin/env bash

# build project for darwin (Mac OS X)

. utils.sh

export GOPATH=`pwd`/../../
export GOOS=darwin
export GOARCH=amd64

build
