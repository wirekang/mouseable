package view

import (
	"bytes"
	"fmt"

	"github.com/JamesHovious/w32"
	"github.com/lxn/walk"
	"github.com/mat/besticon/ico"
	"github.com/pkg/errors"
	"github.com/wirekang/vkmap"

	"github.com/wirekang/mouseable/asset"
	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/def"
)

var mainWindow *walk.MainWindow
var config def.Config

func AlertError(msg string) {
	w32.MessageBox(0, msg, "Mouseable", 0)
}

func Run() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() (err error) {
	mainWindow, err = walk.NewMainWindowWithName("Mouseable")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	config = DI.LoadConfig()

	err = tempUI()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	// prevent window flashing when using hot reloading in development
	mainWindow.SetVisible(!cnst.IsDev)

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
	for fnc, keyCode := range config.FunctionKeyCodeMap {
		txt += fnc.Name + ": "
		t := vkmap.Map[keyCode].VK
		if t == "" {
			t = vkmap.Map[keyCode].Description
		}
		txt += " <" + t + "> "
		txt += "\r\n"
	}
	txt += "\r\n\r\nData\r\n\r\n"
	for data, value := range config.DataValueMap {
		txt += data.Name + ": " + fmt.Sprintf("%f", value) + "\n"
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
			DI.SaveConfig(config)
		},
	)

	return
}
