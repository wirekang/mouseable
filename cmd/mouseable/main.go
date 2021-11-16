package main

import (
	"context"
	"flag"
	"os"

	"golang.org/x/sys/windows/svc"

	"github.com/wirekang/mouseable/internal/check"
	"github.com/wirekang/mouseable/internal/config"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/script"
	"github.com/wirekang/mouseable/internal/winsvc"
)

func main() {
	check.MustWindows()

	register := flag.Bool("register", false, "Register and run service")
	unregister := flag.Bool("unregister", false, "Unregister service")
	reload := flag.Bool(
		"reload", false, "Reload config file at  "+config.FilePath,
	)
	open := flag.Bool(
		"open", false, "Open directory where config.json and debug.log exist.",
	)
	run := flag.Bool(
		"run", false,
		"This flag run mouseable in foreground that usually NOT NEEDED",
	)
	flag.Parse()

	defer func() {
		r := recover()
		if r != nil {
			lg.Errorf("panic: %v", r)
			panic(r)
		}
	}()

	if *register {
		err := script.Register()
		if err != nil {
			panic(err)
		}
	}

	if *unregister {
		err := script.Unregister()
		if err != nil {
			panic(err)
		}
	}

	if *reload {
		err := script.Reload()
		if err != nil {
			panic(err)
		}
	}

	if *open {
		err := script.OpenConfigDir()
		if err != nil {
			panic(err)
		}
	}

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *run {
		lg.Logf("Run")
		err := svc.Run(
			script.ServiceName, winsvc.Handler{
				RunFunc: func(ctx context.Context) {
					lg.Logf("Start")
					<-ctx.Done()
					lg.Logf("Stop")
				},
			},
		)

		if err != nil {
			panic(err)
		}
	}
}
