#!/usr/bin/env bash

# build project for darwin (Mac OS X)

. utils.sh


export GO111MODULE=on
export GOOS=darwin
export GOARCH=amd64

build
