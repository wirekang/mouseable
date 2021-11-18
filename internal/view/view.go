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
		err = errors.Wrap(err, "NewMainWindow")
		return
	}

	if err := mainWindow.SetLayout(walk.NewVBoxLayout()); err != nil {
		panic(err)
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
		panic(err)
	}

	icon, err := walk.NewIconFromImageForDPI(img, 96)
	if err != nil {
		panic(err)
	}

	defer icon.Dispose()

	notifyIcon, err := walk.NewNotifyIcon(mainWindow)
	if err != nil {
		panic(err)
	}

	defer notifyIcon.Dispose()

	if err := notifyIcon.SetIcon(icon); err != nil {
		panic(err)
	}

	if err := notifyIcon.SetToolTip("ToolTip"); err != nil {
		panic(err)
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
	if err := exitAction.SetText("E&xit"); err != nil {
		panic(err)
	}

	exitAction.Triggered().Attach(
		func() { walk.App().Exit(0) },
	)

	if err := notifyIcon.ContextMenu().Actions().Add(exitAction); err != nil {
		panic(err)
	}

	if err := notifyIcon.SetVisible(true); err != nil {
		panic(err)
	}

	if err := notifyIcon.ShowInfo("title", "info"); err != nil {
		panic(err)
	}

	mainWindow.Run()
	return
}
