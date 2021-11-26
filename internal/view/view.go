package view

import (
	"bytes"
	"fmt"
	"os/exec"
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
var config def.Config

func AlertError(msg string) {
	w32.MessageBox(0, msg, "Mouseable", 0)
}

func Run() {
	mainWindow, _ = walk.NewMainWindowWithName("Mouseable " + cnst.VERSION)
	defaultStyle := win.GetWindowLong(mainWindow.Handle(), win.GWL_STYLE)
	win.SetWindowLong(
		mainWindow.Handle(), win.GWL_STYLE, defaultStyle^win.WS_THICKFRAME,
	)
	config = DI.LoadConfig()
	ui()
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

func ui() {
	mainLayout := walk.NewVBoxLayout()
	mainWindow.SetLayout(mainLayout)
	mainWindow.SetWidth(700)
	mainWindow.SetHeight(900)
	uiInfoGroup()
	hotKeyGroup := newVGroupBox(mainWindow, "HotKey")
	newRowText(
		hotKeyGroup,
		"※ Using Win key alone is not recommended. ※\n※ Click left to set key, right to unset.※",
	)
	for i := range def.HotKeyDefs {
		uiHotkey(hotKeyGroup, def.HotKeyDefs[i])
	}

	functionGroup := newVGroupBox(mainWindow, "Function")
	for i := range def.FunctionDefs {
		uiFunction(functionGroup, def.FunctionDefs[i])
	}

	dataGroup := newVGroupBox(mainWindow, "Data")
	for i := range def.DataDefs {
		uiData(dataGroup, def.DataDefs[i])
	}
}

func uiInfoGroup() {
	infoGroup := newVGroupBox(mainWindow, "Info")
	hs, _ := walk.NewHSplitter(infoGroup)
	tl, _ := walk.NewTextLabel(hs)
	tl.SetText(fmt.Sprintf("Version: %s", cnst.VERSION))
	btnGithub, _ := walk.NewPushButton(hs)
	btnGithub.SetText("If you find any bug or have a good idea, visit the GitHub.")

	btnGithub.Clicked().Attach(
		func() {
			exec.Command(
				"rundll32", "url.dll,FileProtocolHandler",
				"https://github.com/wirekang/mouseable",
			).Start()
		},
	)
}

func newVGroupBox(p walk.Container, title string) (gb *walk.GroupBox) {
	gb, _ = walk.NewGroupBox(p)
	gb.SetLayout(walk.NewVBoxLayout())
	gb.SetTitle(title)
	return
}

func newRowText(p walk.Container, txt string) (tl *walk.TextLabel) {
	v, _ := walk.NewVSplitter(p)
	tl, _ = walk.NewTextLabel(v)
	tl.SetText(txt)
	return
}

func uiHotkey(c walk.Container, hkd *def.HotKeyDef) {
	s, _ := walk.NewHSplitter(c)
	name, _ := walk.NewTextLabel(s)
	name.SetText(hkd.Name)
	description, _ := walk.NewTextLabel(s)
	description.SetText(hkd.Description)
	uiCheckbox(
		config.HotKeyMap[hkd].IsControl, "Ctrl", s, func(b bool) {
			t := config.HotKeyMap[hkd]
			t.IsControl = b
			config.HotKeyMap[hkd] = t
			DI.SaveConfig(config)
		},
	)
	uiCheckbox(
		config.HotKeyMap[hkd].IsShift, "Shift", s, func(b bool) {
			t := config.HotKeyMap[hkd]
			t.IsShift = b
			config.HotKeyMap[hkd] = t
			DI.SaveConfig(config)
		},
	)
	uiCheckbox(
		config.HotKeyMap[hkd].IsAlt, "Alt", s, func(b bool) {
			t := config.HotKeyMap[hkd]
			t.IsAlt = b
			config.HotKeyMap[hkd] = t
			DI.SaveConfig(config)
		},
	)
	uiCheckbox(
		config.HotKeyMap[hkd].IsWin, "Win", s, func(b bool) {
			t := config.HotKeyMap[hkd]
			t.IsWin = b
			config.HotKeyMap[hkd] = t
			DI.SaveConfig(config)
		},
	)

	uiKeyCode(
		config.HotKeyMap[hkd].KeyCode, s, func(v uint32) {
			t := config.HotKeyMap[hkd]
			t.KeyCode = v
			config.HotKeyMap[hkd] = t
			DI.SaveConfig(config)
		},
	)
}

func uiFunction(c walk.Container, fnc *def.FunctionDef) {
	s, _ := walk.NewHSplitter(c)
	name, _ := walk.NewTextLabel(s)
	name.SetText(fnc.Name)
	description, _ := walk.NewTextLabel(s)
	description.SetText(fnc.Description)
	uiKeyCode(
		config.FunctionKeyCodeMap[fnc], s, func(u uint32) {
			config.FunctionKeyCodeMap[fnc] = u
			DI.SaveConfig(config)
		},
	)
}

func uiData(c walk.Container, dd *def.DataDef) {
	s, _ := walk.NewHSplitter(c)
	name, _ := walk.NewTextLabel(s)
	name.SetText(dd.Name)
	description, _ := walk.NewTextLabel(s)
	description.SetText(dd.Description)
	uiFloat(
		config.DataValueMap[dd], s, func(f float64) {
			config.DataValueMap[dd] = f
			DI.SaveConfig(config)
		},
	)
}

func uiCheckbox(checked bool, label string, c walk.Container, f func(bool)) {
	cb, _ := walk.NewCheckBox(c)
	cb.SetChecked(checked)
	cb.SetText(label)

	cb.CheckedChanged().Attach(
		func() {
			f(cb.Checked())
		},
	)
}

func uiKeyCode(kc uint32, c walk.Container, cb func(uint32)) {
	btn, _ := walk.NewPushButton(c)
	btn.SetText(getKeyCodeText(kc))
	setF := func(u uint32) {
		btn.SetText(getKeyCodeText(u))
		cb(u)
	}
	btn.MouseUp().Attach(
		func(x, y int, button walk.MouseButton) {
			if button == walk.RightButton {
				setF(0)
			}

			if button == walk.LeftButton {
				openKeyCodeWindow(setF)
			}
		},
	)
}

func uiFloat(v float64, c walk.Container, cb func(float64)) {
	edit, _ := walk.NewTextEdit(c)
	edit.SetText(fmt.Sprintf("%.1f", v))
	btn, _ := walk.NewPushButton(c)
	btn.SetText("Apply")
	btn.Clicked().Attach(
		func() {
			f, err := strconv.ParseFloat(edit.Text(), 64)
			if err != nil {
				return
			}

			edit.SetText(fmt.Sprintf("%.1f", f))
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

var nextClassID = 0

func openKeyCodeWindow(cb func(uint32)) {
	nextClassID++
	className, _ := syscall.UTF16PtrFromString(
		fmt.Sprintf(
			"MouseableKeyCodePopUpWindow%d", nextClassID,
		),
	)
	windowName, _ := syscall.UTF16PtrFromString("Enter Key")
	proc := func(
		hwnd w32.HWND, msg uint32, wparam w32.WPARAM, lparam w32.LPARAM,
	) w32.LRESULT {
		switch msg {
		case w32.WM_ACTIVATE:
			switch wparam {
			case w32.WA_INACTIVE:
				w32.DestroyWindow(hwnd)
				return 0
			}
		case w32.WM_KEYDOWN:
			fmt.Println("KEYDOWN: ", hwnd, wparam)
			cb(uint32(wparam))
			w32.DestroyWindow(hwnd)
			return 0
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
	w32.ShowWindow(hwnd, w32.SW_SHOW)
	lg.Logf("HWND %d\n", hwnd)
}
