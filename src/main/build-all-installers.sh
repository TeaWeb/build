#!/usr/bin/env bash

# build project for all platforms

. utils.sh

export GOPATH=`pwd`/../../

rm -rf ${GOPATH}/installers/*

export GOOS=linux
export GOARCH=386
buildAgentInstaller

export GOOS=linux
export GOARCH=amd64
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
