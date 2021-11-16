package main

import (
	"flag"

	"github.com/wirekang/mouseable/internal/check"
	"github.com/wirekang/mouseable/internal/config"
)

func main() {
	check.MustWindows()

	configPath := flag.String(
		"config", config.DefaultConfigPath, "config file path",
	)
	flag.Parse()

	_, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}
}
