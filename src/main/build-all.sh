#!/usr/bin/env bash

# build project for all platforms

. utils.sh

if [ "$1" = "" ]
then
	. build-all-agents.sh
	. build-all-installers.sh

	. build-linux-32.sh
	. build-linux-64.sh
	. build-linux-arm-64.sh
	. build-linux-arm.sh
	. build-darwin.sh
	. build-windows-32.sh
	. build-windows-64.sh
elif [ "$2" != "" ]
then
		export GOOS=$1
		export GOARCH=$2
		build
elif [ "$2" = "" ]
then
	echo "Usage: ${0} OS ARCH"
fi