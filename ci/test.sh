#!/usr/bin/env bash

set -e

export GOTOOLCHAIN=go1.27.0+auto
export GO111MODULE="on"
go test -race -v ./cmd ./selenoid -coverprofile=coverage.txt -covermode=atomic -coverpkg ./cmd,./selenoid

GOTOOLCHAIN=go1.27.0 go run golang.org/x/vuln/cmd/govulncheck@v1.5.0 ./...
