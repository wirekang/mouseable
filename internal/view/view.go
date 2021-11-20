package view

import (
	"bytes"
	"os"

	"github.com/lxn/walk"
	"github.com/mat/besticon/ico"
	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/asset"
)

var mainWindow *walk.MainWindow

func Run() (err error) {
	mainWindow, err = walk.NewMainWindowWithName("Mouseable")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = initMainWindowLayout()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// prevent window flashing when using hot reload in development
	//
	// checking -dev.exe instead of -dev is due to bug of air
	// https://github.com/cosmtrek/air/issues/207
	if !(len(os.Args) == 2 && os.Args[1] == "-dev.exe") {
		mainWindow.SetVisible(true)
	}

	mainWindow.Closing().Attach(
		func(canceled *bool, reason walk.CloseReason) {
			*canceled = true
			mainWindow.SetVisible(false)
		},
	)

	br := bytes.NewReader(asset.Icon)
	img, err := ico.Decode(br)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	icon, err := walk.NewIconFromImageForDPI(img, 96)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	defer icon.Dispose()

	notifyIcon, err := walk.NewNotifyIcon(mainWindow)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	defer notifyIcon.Dispose()

	err = notifyIcon.SetIcon(icon)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = notifyIcon.SetToolTip("ToolTip")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	notifyIcon.MouseDown().Attach(
		func(x, y int, button walk.MouseButton) {
			if button != walk.LeftButton {
				return
			}

			mainWindow.SetVisible(true)
		},
	)

	exitAction := walk.NewAction()
	err = exitAction.SetText("E&xit")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	exitAction.Triggered().Attach(
		func() { walk.App().Exit(0) },
	)

	err = notifyIcon.ContextMenu().Actions().Add(exitAction)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = notifyIcon.SetVisible(true)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = notifyIcon.ShowInfo("title", "info")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	mainWindow.Run()
	return
}

func initMainWindowLayout() (err error) {
	err = mainWindow.SetLayout(walk.NewVBoxLayout())
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
