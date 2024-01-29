#!/bin/sh

#
# used to spawn a new redis server
set -e
tmpFile=$(mktemp)
go build -o "$tmpFile" main.go
exec "$tmpFile" "$@"
