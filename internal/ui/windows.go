//go:build windows

package ui

import (
	"github.com/JamesHovious/w32"
	"github.com/lxn/walk"
	"github.com/mat/besticon/ico"

	"github.com/wirekang/mouseable/internal/cnst"
)

func (m *manager) StartBackground() {
	m.runNotifyIcon()
}

func showAlert(msg string) {
	w32.MessageBox(0, msg, "Mouseable", 0)
}

func showError(msg string) {
	w32.MessageBox(0, msg+"\n\nPress Ctrl+C to copy message.", "Mouseable", w32.MB_ICONERROR)
}

func (m *manager) runNotifyIcon() {
	br, err := cnst.FrontFS.Open("favicon.ico")
	if err != nil {
		panic(err)
	}

	img, err := ico.Decode(br)
	if err != nil {
		panic(err)
	}

	icon, err := walk.NewIconFromImageForDPI(img, 96)
	if err != nil {
		panic(err)
	}

	defer icon.Dispose()

	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		panic(err)
	}

	notifyIcon, err := walk.NewNotifyIcon(mainWindow)
	if err != nil {
		panic(err)
	}

	defer notifyIcon.Dispose()

	err = notifyIcon.SetIcon(icon)
	if err != nil {
		panic(err)
	}

	notifyIcon.MouseDown().Attach(
		func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				go m.Open()
			}
		},
	)

	err = notifyIcon.SetVisible(true)
	if err != nil {
		panic(err)
	}
	mainWindow.Run()
}
