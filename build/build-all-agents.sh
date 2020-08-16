#!/usr/bin/env bash

# build project for all platforms

CWD=$(dirname $0)
ROOT=$CWD/../

. $CWD/utils.sh

export GOPATH=""
export GO111MODULE=on

if [ "$1" = "" ]
then
	rm -rf ${ROOT}/web/upgrade/*

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