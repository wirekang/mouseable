package view

import (
	"bytes"
	"io/fs"
	"net"
	"net/http"

	"github.com/JamesHovious/w32"
	"github.com/lxn/walk"
	"github.com/mat/besticon/ico"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/def"
)

var config def.Config
var host string

func AlertError(msg string) {
	w32.MessageBox(0, msg, "Mouseable", 0)
}

func serve() {
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

func Run() {
	mustChrome()
	go serve()
	config = DI.LoadConfig()

	defer func() {

	}()
	runNotifyIcon()
}

func runNotifyIcon() {
	bs, err := cnst.AssetFS.ReadFile("assets/favicon.ico")
	if err != nil {
		panic(err)
	}

	br := bytes.NewReader(bs)
	img, _ := ico.Decode(br)
	icon, _ := walk.NewIconFromImageForDPI(img, 96)
	defer icon.Dispose()

	mainWindow, _ := walk.NewMainWindow()

	notifyIcon, _ := walk.NewNotifyIcon(mainWindow)
	defer notifyIcon.Dispose()

	notifyIcon.SetIcon(icon)
	notifyIcon.SetToolTip("ToolTip")

	notifyIcon.MouseDown().Attach(
		func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				go open()
			}
		},
	)

	exitAction := walk.NewAction()
	exitAction.SetText("Exit")

	exitAction.Triggered().Attach(
		func() {
			ui.Close()
			walk.App().Exit(0)
		},
	)

	notifyIcon.ContextMenu().Actions().Add(exitAction)
	err = notifyIcon.SetVisible(true)
	if err != nil {
		panic(err)
	}
	mainWindow.Run()
}
