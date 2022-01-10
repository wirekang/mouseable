package main

import (
	"embed"
	"os"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/logic"
)

//go:embed assets
var Asset embed.FS
var VERSION = "x.x.x"

func main() {
	cnst.VERSION = VERSION
	cnst.AssetFS = Asset

	// checking -dev.exe instead of -dev is due to bug of air.
	// https://github.com/cosmtrek/air/issues/207
	if len(os.Args) == 2 && (os.Args[1] == "-dev.exe" || os.Args[1] == "-dev") {
		cnst.IsDev = true
	}

	logic.Run()
}
