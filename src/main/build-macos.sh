#!/usr/bin/env bash

. utils.sh

export GOPATH=`pwd`/../../
export GOOS=darwin
export GOARCH=amd64

build
