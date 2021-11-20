package view

import (
	"bytes"
	"fmt"

	"github.com/lxn/walk"
	"github.com/mat/besticon/ico"
	"github.com/pkg/errors"
	"github.com/wirekang/vkmap"

	"github.com/wirekang/mouseable/asset"
	"github.com/wirekang/mouseable/internal/lg"
)

var mainWindow *walk.MainWindow
var keymap map[string][]uint32
var data map[string]string

func Run() (err error) {
	mainWindow, err = walk.NewMainWindowWithName("Mouseable")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	keymap, data, err = DI.LoadData()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = tempUI()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	// prevent window flashing when using hot reloading in development
	mainWindow.SetVisible(!lg.IsDev)

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

func tempUI() (err error) {
	vbox := walk.NewVBoxLayout()
	err = mainWindow.SetLayout(vbox)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	txt := "\r\nKeymap\r\n\r\n"
	for name, keycodes := range keymap {
		txt += name + ": "
		for _, kc := range keycodes {
			t := vkmap.Map[kc].VK
			if t == "" {
				t = vkmap.Map[kc].Description
			}
			txt += " <" + t + "> "
		}
		txt += "\r\n"
	}
	txt += "\r\n\r\nData\r\n\r\n"
	for key, value := range data {
		txt += key + ": " + value + "\n"
	}

	label, err := walk.NewTextLabel(vbox.Container())
	err = label.SetText(txt)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	btn, err := walk.NewPushButton(vbox.Container())
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = btn.SetText("Save to appdata")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	btn.MouseDown().Attach(
		func(_, _ int, _ walk.MouseButton) {
			err := DI.SaveData(keymap, data)
			if err != nil {
				fmt.Println(err)
				label.SetText(err.Error())
			}
		},
	)

	return
}
