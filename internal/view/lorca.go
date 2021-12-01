package view

import (
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/zserge/lorca"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/lg"
)

var isOpened bool
var host string
var ui lorca.UI

func mustChrome() {
	if lorca.LocateChrome() == "" {
		AlertError("Chromium browser not found. Mouseable can't render GUI. Please install Chrome or Edge.")
		os.Exit(1)
	}
}

func serveAndOpen() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	host = ln.Addr().String()
	go open()
	defer ln.Close()
	sub, err := fs.Sub(cnst.AssetFS, "assets")
	if err != nil {
		panic(err)
	}

	http.Serve(ln, http.FileServer(http.FS(sub)))
}

func open() {
	if isOpened {
		return
	}

	url := "http://" + host
	lg.Logf("Open: %s", url)

	var err error
	ui, err = lorca.New(url, "", 800, 800, "--disable-features=Translate")
	if err != nil {
		panic(err)
	}

	isOpened = true
	defer func() {
		isOpened = false
		ui.Close()
		lg.Logf("Close")
	}()

	err = ui.Bind("__loadBind__", __loadBind__)
	if err != nil {
		panic(err)
	}

	err = ui.Bind("__getKeyCode__", __getKeyCode__)
	if err != nil {
		panic(err)
	}

	err = ui.Bind("__changeFunction__", __changeFunction__)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case <-sigChan:
		exit()
	case <-ui.Done():
	}
	return
}
