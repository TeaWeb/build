#!/usr/bin/env bash

# build project for all platforms

. utils.sh

export GOPATH=`pwd`/../../
export GO111MODULE=off

if [ "$1" = "" ]
then
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

	export GOOS=linux
	export GOARCH=mips64
	buildAgent

	export GOOS=linux
	export GOARCH=mips64le
	buildAgent
elif [ "$2" != "" ]
then
		export GOOS=$1
		export GOARCH=$2
		buildAgent
elif [ "$2" = "" ]
then
	echo "Usage: ${0} OS ARCH"
fi