package view

import (
	"bytes"
	"os"

	"github.com/JamesHovious/w32"
	"github.com/lxn/walk"
	"github.com/mat/besticon/ico"

	"github.com/wirekang/mouseable/internal/cnst"
)

func AlertError(msg string) {
	w32.MessageBox(0, msg, "Mouseable", 0)
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

	exitAction.Triggered().Attach(exit)

	notifyIcon.ContextMenu().Actions().Add(exitAction)
	err = notifyIcon.SetVisible(true)
	if err != nil {
		panic(err)
	}
	mainWindow.Run()
}

func exit() {
	ui.Close()
	os.Exit(0)
}
