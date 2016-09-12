#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

mkdir -p /home/etcinit/go/src/github.com/etcinit
ln -s $(pwd) /home/drydock/go/src/github.com/etcinit/
cd /home/etcinit/go/src/github.com/etcinit/phabulous

go get ./...
go build cmd/phabulous/phabulous.go
