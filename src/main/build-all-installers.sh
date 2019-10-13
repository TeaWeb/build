#!/usr/bin/env bash

# build project for all platforms

. utils.sh

export GOPATH=`pwd`/../../
export GO111MODULE=off

rm -rf ${GOPATH}/installers/*

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
