package main

import (
	"flag"
	"os"

	"github.com/wirekang/winsvc/internal/lg"
	"github.com/wirekang/winsvc/internal/must"
	"github.com/wirekang/winsvc/internal/script"
)

func main() {
	must.Windows()

	openConfig := flag.Bool(
		"config", false, "Open config file.",
	)
	flag.Parse()

	defer func() {
		r := recover()
		if r != nil {
			lg.Errorf("panic: %v", r)
			panic(r)
		}
	}()

	if *openConfig {
		err := script.OpenConfigFile()
		if err != nil {
			panic(err)
		}
	}

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

}
