#!/usr/bin/env bash

# build project for all platforms

. utils.sh

. build-all-agents.sh
. build-all-installers.sh

. build-linux-32.sh
. build-linux-64.sh
. build-darwin.sh
. build-windows-32.sh
. build-windows-64.sh