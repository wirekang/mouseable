package view

import (
	"bytes"
	"image/png"

	"github.com/lxn/walk"
	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/asset"
)

func Init() (err error) {
	mainWindow, err := walk.NewMainWindowWithName("Mouseable")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = mainWindow.SetLayout(walk.NewVBoxLayout())
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	mainWindow.SetVisible(true)
	mainWindow.Closing().Attach(
		func(canceled *bool, reason walk.CloseReason) {
			*canceled = true
			mainWindow.SetVisible(false)
		},
	)

	br := bytes.NewReader(asset.Icon)
	img, err := png.Decode(br)
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
