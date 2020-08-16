#!/usr/bin/env bash

# build project for all platforms

CWD=$(dirname $0)
ROOT=$CWD/../

. $CWD/utils.sh

if [ "$1" = "" ]
then
	export GO111MODULE=on

	. $CWD/build-all-agents.sh
	. $CWD/build-all-installers.sh

	GOOS=linux GOARCH=386 build
	GOOS=linux GOARCH=amd64 build
	GOOS=linux GOARCH=arm build
	GOOS=linux GOARCH=arm64 build
	GOOS=linux GOARCH=mips64 build
	GOOS=linux GOARCH=mips64le build
	GOOS=darwin GOARCH=amd64 build
	GOOS=windows GOARCH=386 build
	GOOS=windows GOARCH=amd64 build
elif [ "$2" != "" ]
then
		export GOOS=$1
		export GOARCH=$2
		export GO111MODULE=on
		build
elif [ "$2" = "" ]
then
	echo "Usage: ${0} OS ARCH"
fi