package main

import (
	"embed"
	"io/fs"
	"os"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/logic"
)

//go:embed assets
var Asset embed.FS
var VERSION = "x.x.x"

func main() {

	cnst.VERSION = VERSION
	initFS()

	if len(os.Args) == 2 && os.Args[1] == "-dev" {
		cnst.IsDev = true
	}

	logic.Run()
}

func initFS() {
	var err error
	cnst.AssetFS, err = fs.Sub(Asset, "assets")
	if err != nil {
		panic(err)
	}

	cnst.FrontFS, err = fs.Sub(cnst.AssetFS, "front")
	if err != nil {
		panic(err)
	}

	cnst.DefaultConfigsFS, err = fs.Sub(cnst.DefaultConfigsFS, "defaultConfigs")
	if err != nil {
		panic(err)
	}
}
