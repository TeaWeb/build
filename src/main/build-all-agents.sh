#!/usr/bin/env bash

# build project for all platforms

. utils.sh

export GOPATH=`pwd`/../../

rm -rf ${GOPATH}/upgrade/*

export GOOS=linux
export GOARCH=amd64
buildAgent

export GOOS=linux
export GOARCH=386
buildAgent

export GOOS=linux
export GOARCH=arm64
buildAgent

export GOOS=linux
export GOARCH=arm
buildAgent

export GOOS=darwin
export GOARCH=amd64
buildAgent

export GOOS=windows
export GOARCH=386
buildAgent

export GOOS=windows
export GOARCH=amd64
buildAgent

export GOOS=freebsd
export GOARCH=amd64
buildAgent
