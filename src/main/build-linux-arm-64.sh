#!/usr/bin/env bash

# build project for linux 64 bit

. utils.sh

export GOPATH=`pwd`/../../
export GO111MODULE=off
export GOOS=linux
export GOARCH=arm64

build
