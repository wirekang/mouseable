package main

import (
	"flag"
	"os"

	"github.com/wirekang/mouseable/internal/check"
	"github.com/wirekang/mouseable/internal/config"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/svc"
)

func main() {
	check.MustWindows()

	register := flag.Bool("register", false, "Register and run service")
	unregister := flag.Bool("unregister", false, "Unregister service")
	reload := flag.Bool(
		"reload", false, "Reload config file at  "+config.ConfigPath,
	)
	run := flag.Bool(
		"run", false,
		"This flag run mouseable in foreground that usually NOT NEEDED",
	)
	flag.Parse()

	if *register {
		err := svc.Register()
		if err != nil {
			panic(err)
		}
	}

	if *unregister {
		err := svc.Unregister()
		if err != nil {
			panic(err)
		}
	}

	if *reload {
		err := svc.Reload()
		if err != nil {
			panic(err)
		}
	}

	if !(*run || *register || *unregister || *reload) {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *run {
		lg.Logf("start")
		_, err := config.Load()
		if err != nil {
			lg.Errorf("config.Load: %v", err)
			os.Exit(1)
		}
	}
}
