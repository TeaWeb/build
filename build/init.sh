#!/usr/bin/env bash

# initialize project
export GO111MODULE=on
export GOPROXY=direct

# download
go mod tidy

echo "[done]"