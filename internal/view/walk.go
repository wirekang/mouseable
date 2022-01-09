package view

import (
	"bytes"

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

	_ = notifyIcon.SetIcon(icon)

	notifyIcon.MouseDown().Attach(
		func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				go openUI()
			}
		},
	)

	err = notifyIcon.SetVisible(true)
	if err != nil {
		panic(err)
	}
	mainWindow.Run()
}
