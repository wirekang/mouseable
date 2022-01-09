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

var isOpen = false

func mustChrome() {
	if lorca.LocateChrome() == "" {
		AlertError("Chromium browser not found. Mouseable can't render GUI. Please install Chrome or Edge.")
		lorca.PromptDownload()
		os.Exit(1)
	}
}

func openUI() {
	if isOpen {
		lg.Logf("Window is already open.")
		return
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	sub, err := fs.Sub(cnst.AssetFS, "assets")
	if err != nil {
		panic(err)
	}

	go func() {
		err = http.Serve(listener, http.FileServer(http.FS(sub)))
		if err != nil {
			lg.Errorf("http.Serve: %v", err)
		}
	}()

	host := "http://" + listener.Addr().String()
	ui, err := lorca.New(host, "", 800, 800, "--disable-features=Translate")
	if err != nil {
		panic(err)
	}

	err = bindLorca(ui)
	if err != nil {
		panic(err)
	}

	isOpen = true
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case <-sigChan:
		_ = ui.Close()
		os.Exit(0)
	case <-ui.Done():
	}

	err = listener.Close()
	if err != nil {
		lg.Errorf("listener.Close(): %v", err)
	}

	isOpen = false
	err = ui.Close()
	if err != nil {
		lg.Errorf("ui.Close(): %v", err)
	}
	lg.Logf("Close UI")
}
