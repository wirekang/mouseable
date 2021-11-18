#!/bin/bash
cd "$(dirname "$(dirname "$0")}")" || return
mkdir build || ls
rsrc -ico assets/icon.ico -manifest mouseable.manifest -o cmd/mouseable/rsrc.syso
go build -ldflags="-H windowsgui" -o ./build/out.exe ./cmd/mouseable/main.go
