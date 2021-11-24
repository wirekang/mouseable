#!/bin/bash
cd "$(dirname "$(dirname "$0")}")" || return
VERSION="$(cat version)"
mkdir -p build
rsrc -ico asset/icon.ico -manifest mouseable.manifest
go build -ldflags="-H windowsgui -X main.VERSION=$VERSION" -o ./build/portable.exe
