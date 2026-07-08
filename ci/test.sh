#!/usr/bin/env bash

set -e

export GOTOOLCHAIN=go1.26.5+auto
export GO111MODULE="on"
go test -race -v github.com/aerokube/cm/selenoid -coverprofile=coverage.txt -covermode=atomic -coverpkg github.com/aerokube/cm/selenoid

GOTOOLCHAIN=go1.26.5 go run golang.org/x/vuln/cmd/govulncheck@v1.5.0 ./...
