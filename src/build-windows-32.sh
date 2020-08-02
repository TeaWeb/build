#!/usr/bin/env bash

# build project for linux 64 bit

. utils.sh


export GO111MODULE=on
export GOOS=windows
export GOARCH=386

build
