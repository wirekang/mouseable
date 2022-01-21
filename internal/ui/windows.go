//go:build windows

package ui

import (
	"io/fs"

	"github.com/JamesHovious/w32"
	"github.com/lxn/walk"
	"github.com/mat/besticon/ico"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/lg"
)

func (m *manager) Run() {
	m.runNotifyIcon()
}

func showAlert(msg string) {
	w32.MessageBox(0, msg, "Mouseable", 0)
}

func showError(msg string) {
	w32.MessageBox(0, msg+"\n\nPress Ctrl+C to copy message.", "Mouseable", w32.MB_ICONERROR)
}

func (m *manager) setTrayIconEnabled(b bool) {
	if b {
		err := m.notifyIcon.SetIcon(m.iconEnabled)
		if err != nil {
			lg.Errorf("SetIcon: %v", err)
		}
	} else {
		err := m.notifyIcon.SetIcon(m.iconDisabled)
		if err != nil {
			lg.Errorf("SetIcon: %v", err)
		}
	}
}

func (m *manager) runNotifyIcon() {
	m.iconDisabled = decodeIco(cnst.FrontFS, "favicon.ico")
	defer m.iconDisabled.Dispose()

	m.iconEnabled = decodeIco(cnst.AssetFS, "enabled.ico")
	defer m.iconEnabled.Dispose()

	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		panic(err)
	}

	m.notifyIcon, err = walk.NewNotifyIcon(mainWindow)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = m.notifyIcon.Dispose()
	}()

	err = m.notifyIcon.SetIcon(m.iconDisabled)
	if err != nil {
		panic(err)
	}
	m.notifyIcon.MouseDown().Attach(
		func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				go m.Open()
			}
		},
	)

	err = m.notifyIcon.SetVisible(true)
	if err != nil {
		panic(err)
	}

	mainWindow.Run()
}

func decodeIco(fss fs.FS, name string) (i *walk.Icon) {
	br, err := fss.Open(name)
	if err != nil {
		panic(err)
	}

	img, err := ico.Decode(br)
	if err != nil {
		panic(err)
	}

	i, err = walk.NewIconFromImageForDPI(img, 96)
	if err != nil {
		panic(err)
	}

	return
}
