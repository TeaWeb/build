#!/usr/bin/env bash

# build project for all platforms

CWD=$(dirname $0)
ROOT=$CWD/../

. $CWD/utils.sh

export GOPATH=""
export GO111MODULE=on

#rm -rf ${ROOT}/web/installers/*

export GOOS=linux
export GOARCH=386
buildAgentInstaller

export GOOS=linux
export GOARCH=amd64
buildAgentInstaller

export GOOS=linux
export GOARCH=arm64
buildAgentInstaller

export GOOS=linux
export GOARCH=arm
buildAgentInstaller

export GOOS=linux
export GOARCH=mips64
buildAgentInstaller

export GOOS=linux
export GOARCH=mips64le
buildAgentInstaller

export GOOS=darwin
export GOARCH=amd64
buildAgentInstaller

#export GOOS=windows
#export GOARCH=386
#buildAgentInstaller

#export GOOS=windows
#export GOARCH=amd64
#buildAgentInstaller
