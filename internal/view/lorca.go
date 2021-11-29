package view

import (
	"os"

	"github.com/zserge/lorca"

	"github.com/wirekang/mouseable/internal/lg"
)

var isOpened bool
var ui lorca.UI

func mustChrome() {
	if lorca.LocateChrome() == "" {
		AlertError("Google Chrome not found. Mouseable can't render GUI.")
		os.Exit(1)
	}
}

func open() {
	if isOpened {
		return
	}

	url := "http://" + host
	lg.Logf("Open: %s", url)

	var err error
	ui, err = lorca.New(url, "", 800, 800)
	if err != nil {
		panic(err)
	}

	isOpened = true
	defer func() {
		isOpened = false
		ui.Close()
		lg.Logf("Close")
	}()
	<-ui.Done()
	return
}
