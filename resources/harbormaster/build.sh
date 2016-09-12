#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

mkdir -p /home/drydock/go/src/github.com/etcinit
ln -sf $(pwd) /home/drydock/go/src/github.com/etcinit/
cd /home/drydock/go/src/github.com/etcinit/phabulous

go get ./...
go build cmd/phabulous/phabulous.go
go test ./...
