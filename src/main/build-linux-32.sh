#!/usr/bin/env bash

# build project for linux 32 bit

. utils.sh

export GOPATH=`pwd`/../../
export GO111MODULE=off
export GOOS=linux
export GOARCH=386

build
