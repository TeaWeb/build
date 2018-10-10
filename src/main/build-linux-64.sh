#!/usr/bin/env bash

. utils.sh

export GOPATH=`pwd`/../../
export GOOS=linux
export GOARCH=amd64

build
