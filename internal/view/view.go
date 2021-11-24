package view

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/JamesHovious/w32"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"github.com/mat/besticon/ico"
	"github.com/wirekang/vkmap"

	"github.com/wirekang/mouseable/asset"
	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/lg"
)

var mainWindow *walk.MainWindow

func AlertError(msg string) {
	w32.MessageBox(0, msg, "Mouseable", 0)
}

func Run() {
	mainWindow, _ = walk.NewMainWindowWithName("Mouseable " + cnst.VERSION)
	mainWindow.SetWidth(700)
	mainWindow.SetHeight(900)
	defaultStyle := win.GetWindowLong(mainWindow.Handle(), win.GWL_STYLE)
	win.SetWindowLong(
		mainWindow.Handle(), win.GWL_STYLE, defaultStyle^win.WS_THICKFRAME,
	)
	config := DI.LoadConfig()
	tempUI(&config)
	mainWindow.SetVisible(true)

	mainWindow.Closing().Attach(
		func(canceled *bool, reason walk.CloseReason) {
			*canceled = true
			mainWindow.SetVisible(false)
		},
	)

	br := bytes.NewReader(asset.Icon)
	img, _ := ico.Decode(br)
	icon, _ := walk.NewIconFromImageForDPI(img, 96)
	mainWindow.SetIcon(icon)
	defer icon.Dispose()

	notifyIcon, _ := walk.NewNotifyIcon(mainWindow)
	defer notifyIcon.Dispose()

	notifyIcon.SetIcon(icon)
	notifyIcon.SetToolTip("ToolTip")

	notifyIcon.MouseDown().Attach(
		func(x, y int, button walk.MouseButton) {
			if button != walk.LeftButton {
				return
			}

			mainWindow.SetVisible(true)
		},
	)

	exitAction := walk.NewAction()
	exitAction.SetText("E&xit")

	exitAction.Triggered().Attach(
		func() { walk.App().Exit(0) },
	)

	notifyIcon.ContextMenu().Actions().Add(exitAction)
	notifyIcon.SetVisible(true)
	notifyIcon.ShowInfo("title", "info")

	mainWindow.Run()
}

func tempUI(config *def.Config) {
	mainWindow.SetLayout(walk.NewVBoxLayout())
	hotKeyGroup, _ := walk.NewGroupBox(mainWindow)
	hotKeyGroup.SetLayout(walk.NewVBoxLayout())
	hotKeyGroup.SetTitle("HotKey")
	helpRow, _ := walk.NewVSplitter(hotKeyGroup)
	helpText, _ := walk.NewTextLabel(helpRow)
	helpText.SetText("※ Using Win key alone is not recommended ※\n※ Click left to set key, right to unset.")
	var hkDefs []*def.HotKeyDef
	for _, hkd := range def.HotKeyNameMap {
		hkDefs = append(hkDefs, hkd)
	}
	sort.Slice(
		hkDefs, func(i, j int) bool {
			return hkDefs[i].Order < hkDefs[j].Order
		},
	)
	for _, hkd := range hkDefs {
		hotKey(config, hotKeyGroup, hkd)
	}

	functionGroup, _ := walk.NewGroupBox(mainWindow)
	functionGroup.SetLayout(walk.NewVBoxLayout())
	functionGroup.SetTitle("Function")
	var fncDefs []*def.FunctionDef
	for _, fncDef := range def.FunctionNameMap {
		fncDefs = append(fncDefs, fncDef)
	}
	sort.Slice(
		fncDefs, func(i, j int) bool {
			return fncDefs[i].Order < fncDefs[j].Order
		},
	)
	for _, fncDef := range fncDefs {
		function(config, functionGroup, fncDef)
	}

	dataGroup, _ := walk.NewGroupBox(mainWindow)
	dataGroup.SetLayout(walk.NewVBoxLayout())
	dataGroup.SetTitle("Data")
	var dataDefs []*def.DataDef
	for _, dataDef := range def.DataNameMap {
		dataDefs = append(dataDefs, dataDef)
	}
	sort.Slice(
		dataDefs, func(i, j int) bool {
			return dataDefs[i].Order < dataDefs[j].Order
		},
	)
	for _, dataDef := range dataDefs {
		data(config, dataGroup, dataDef)
	}

	btn, _ := walk.NewPushButton(mainWindow)
	btn.SetText("Save")

	btn.MouseDown().Attach(
		func(_, _ int, _ walk.MouseButton) {
			DI.SaveConfig(*config)
		},
	)

	return
}

func hotKey(config *def.Config, c walk.Container, hkd *def.HotKeyDef) {
	s, _ := walk.NewHSplitter(c)
	name, _ := walk.NewTextLabel(s)
	name.SetText(hkd.Name)
	description, _ := walk.NewTextLabel(s)
	description.SetText(hkd.Description)
	checkBox(
		config.HotKeyMap[hkd].IsControl, "Ctrl", s, func(b bool) {
			t := config.HotKeyMap[hkd]
			t.IsControl = b
			config.HotKeyMap[hkd] = t
		},
	)
	checkBox(
		config.HotKeyMap[hkd].IsShift, "Shift", s, func(b bool) {
			t := config.HotKeyMap[hkd]
			t.IsShift = b
			config.HotKeyMap[hkd] = t
		},
	)
	checkBox(
		config.HotKeyMap[hkd].IsAlt, "Alt", s, func(b bool) {
			t := config.HotKeyMap[hkd]
			t.IsAlt = b
			config.HotKeyMap[hkd] = t
		},
	)
	checkBox(
		config.HotKeyMap[hkd].IsWin, "Win", s, func(b bool) {
			t := config.HotKeyMap[hkd]
			t.IsWin = b
			config.HotKeyMap[hkd] = t
		},
	)

	keycode(
		config.HotKeyMap[hkd].KeyCode, s, func(v uint32) {
			t := config.HotKeyMap[hkd]
			t.KeyCode = v
			config.HotKeyMap[hkd] = t
		},
	)

}

func function(config *def.Config, c walk.Container, fnc *def.FunctionDef) {
	s, _ := walk.NewHSplitter(c)
	name, _ := walk.NewTextLabel(s)
	name.SetText(fnc.Name)
	description, _ := walk.NewTextLabel(s)
	description.SetText(fnc.Description)
	keycode(
		config.FunctionKeyCodeMap[fnc], s, func(u uint32) {
			config.FunctionKeyCodeMap[fnc] = u
		},
	)
}

func data(config *def.Config, c walk.Container, dd *def.DataDef) {
	s, _ := walk.NewHSplitter(c)
	name, _ := walk.NewTextLabel(s)
	name.SetText(dd.Name)
	description, _ := walk.NewTextLabel(s)
	description.SetText(dd.Description)
	float(
		config.DataValueMap[dd], s, func(f float64) {
			config.DataValueMap[dd] = f
		},
	)
}

func checkBox(checked bool, label string, c walk.Container, f func(bool)) {
	cb, _ := walk.NewCheckBox(c)
	cb.SetChecked(checked)
	cb.SetText(label)

	cb.CheckedChanged().Attach(
		func() {
			f(cb.Checked())
		},
	)
}

func keycode(kc uint32, c walk.Container, cb func(uint32)) {
	btn, _ := walk.NewPushButton(c)
	btn.SetText(getKeyCodeText(kc))
	btn.MouseUp().Attach(
		func(x, y int, button walk.MouseButton) {
			if button == walk.RightButton {
				btn.SetText(getKeyCodeText(0))
				cb(0)
			}

			if button == walk.LeftButton {
				openKeyCodeWindow(
					func(u uint32) {
						btn.SetText(getKeyCodeText(u))
						cb(u)
					},
				)
			}
		},
	)
	btn.Clicked().Attach(
		func() {

		},
	)
}

func float(v float64, c walk.Container, cb func(float64)) {
	edit, _ := walk.NewTextEdit(c)
	edit.SetText(fmt.Sprintf("%.2f", v))
	btn, _ := walk.NewPushButton(c)
	btn.SetText("Apply")
	btn.Clicked().Attach(
		func() {
			f, err := strconv.ParseFloat(edit.Text(), 64)
			if err != nil {
				return
			}

			edit.SetText(fmt.Sprintf("%.2f", f))
			cb(f)
		},
	)

}

func getKeyCodeText(kc uint32) (txt string) {
	if kc == 0 {
		return "-"
	}
	d, ok := vkmap.Map[kc]
	txt = d.VK
	if d.VK == "" {
		txt = d.Description
	}
	if !ok {
		txt = "wrong keyCode"
	}
	return
}

func openKeyCodeWindow(cb func(uint32)) {
	className, _ := syscall.UTF16PtrFromString("MouseableKeyCodePopUpWindow")
	windowName, _ := syscall.UTF16PtrFromString("Enter Key")
	proc := func(
		hwnd w32.HWND, msg uint32, wparam w32.WPARAM, lparam w32.LPARAM,
	) w32.LRESULT {
		switch msg {
		case w32.WM_ACTIVATE:
			switch wparam {
			case w32.WA_INACTIVE:
				w32.DestroyWindow(hwnd)
			}
		case w32.WM_KEYDOWN:
			cb(uint32(wparam))
			w32.DestroyWindow(hwnd)
		}
		return w32.LRESULT(
			w32.DefWindowProc(
				hwnd, msg, uintptr(wparam),
				uintptr(lparam),
			),
		)
	}

	var class w32.WNDCLASSEX
	class.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	class.Background = w32.COLOR_BACKGROUND
	class.ClassName = className
	class.Size = uint32(unsafe.Sizeof(class))
	class.WndProc = syscall.NewCallback(proc)
	w32.RegisterClassEx(&class)
	cx, cy, _ := w32.GetCursorPos()
	w := 300
	h := 40
	hwnd := w32.CreateWindowEx(
		w32.WS_EX_OVERLAPPEDWINDOW,
		className,
		windowName,
		w32.WS_POPUP|w32.WS_CAPTION,
		cx-w/2, cy+10, w, h,
		w32.HWND(mainWindow.Handle()), 0, 0, nil,
	)
	w32.ShowWindow(hwnd, w32.SW_SHOWNORMAL)
	lg.Logf("HWND %d\n", hwnd)
}
