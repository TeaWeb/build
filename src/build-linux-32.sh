#!/usr/bin/env bash

# build project for linux 32 bit

. utils.sh


export GO111MODULE=on
export GOOS=linux
export GOARCH=386

build
