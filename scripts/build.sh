#!/bin/bash
cd "$(dirname "$(dirname "$0")}")" || return
mkdir -p build
rsrc -ico asset/icon.ico -manifest mouseable.manifest
go build -ldflags="-H windowsgui" -o ./build/out.exe
