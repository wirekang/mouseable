//go:build windows

package hook

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/JamesHovious/w32"
	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/lg"
)

func New() di.HookManager {
	return &manager{}
}

type manager struct {
	hHookKeyboard, hHookMouse w32.HHOOK
	onKeyListener             func(di.HookKeyInfo) bool
	onCursorMoveListener      func(info di.Point)
}

func (m *manager) SetOnKeyListener(f func(di.HookKeyInfo) bool) {
	m.onKeyListener = f
}

func (m *manager) SetOnCursorMoveListener(f func(di.Point)) {
	m.onCursorMoveListener = f
}

func (m *manager) SetCursorPosition(x, y int) {
	w32.SetCursorPos(x, y)
}

func (m *manager) AddCursorPosition(dx, dy int) {
	if dx == 0 && dy == 0 {
		return
	}

	sendMouseInput(int32(dx), int32(dy), 0, w32.MOUSEEVENTF_MOVE)
}

func (m *manager) CursorPosition() (x, y int) {
	x, y, _ = w32.GetCursorPos()
	return
}

func (m *manager) MouseDown(button di.MouseButton) {
	sendMouseInput(0, 0, 0, getMouseDownFlag(button))
}

func (m *manager) MouseUp(button di.MouseButton) {
	sendMouseInput(0, 0, 0, getMouseUpFlag(button))
}

func (m *manager) MouseWheel(amount int, isHorizontal bool) {
	if amount != 0 {
		sendMouseInput(0, 0, uint32(amount), getMouseWheelFlag(isHorizontal))
	}
}

func (m *manager) Install() {
	hMod, err := syscall.LoadLibrary("user32.dll")
	if err != nil {
		err = errors.WithStack(err)
		panic(err)
	}

	m.hHookKeyboard = w32.SetWindowsHookEx(w32.WH_KEYBOARD_LL, m.keyboardProc, w32.HINSTANCE(hMod), 0)
	if m.hHookKeyboard == 0 {
		panic(fmt.Sprintf("Keyboard hook failed"))
	}

	lg.Printf("KeyboardHook: %v", m.hHookKeyboard)

	m.hHookMouse = w32.SetWindowsHookEx(w32.WH_MOUSE_LL, m.mouseProc, w32.HINSTANCE(hMod), 0)
	if m.hHookMouse == 0 {
		panic(fmt.Sprintf("Mouse hook failed"))
	}

	lg.Printf("MouseHook: %v", m.hHookMouse)
	return
}

func (m *manager) Uninstall() {
	ok := w32.UnhookWindowsHookEx(m.hHookKeyboard)
	if !ok {
		lg.Errorf("Unhook keyboard failed.")
	}
	ok = w32.UnhookWindowsHookEx(m.hHookMouse)
	if !ok {
		lg.Errorf("Unhook mouse failed.")
	}
}

func (m *manager) keyboardProc(code int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	start := time.Now().UnixMilli()
	defer func() {
		dur := time.Now().UnixMilli() - start
		if dur > 100 {
			lg.Errorf("Hook too long: %dms", dur)
		}
	}()

	data := *(*w32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	// fExtended := data.Flags & (w32.KF_EXTENDED >> 8)
	// fLowerInjected := data.Flags & 0x00000002
	// fInjected := data.Flags & 0x00000010
	// fAltDown := data.Flags & (w32.KF_ALTDOWN >> 8)
	fUp := data.Flags & (w32.KF_UP >> 8)

	isDown := fUp == 0
	txt := getKey(uint32(data.ScanCode))
	if m.onKeyListener(
		di.HookKeyInfo{
			Key:    txt,
			IsDown: isDown,
		},
	) {
		return 1
	}
	return w32.CallNextHookEx(0, code, wParam, lParam)
}

func (m *manager) mouseProc(code int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	data := *(*w32.MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	go m.onCursorMoveListener(
		di.Point{
			X: int(data.Pt.X),
			Y: int(data.Pt.Y),
		},
	)
	return w32.CallNextHookEx(0, code, wParam, lParam)
}

func sendMouseInput(dx, dy int32, mouseData uint32, flags ...uint32) {
	var dwFlags uint32
	for _, f := range flags {
		dwFlags |= f
	}

	input := []w32.INPUT{
		{
			Type: w32.INPUT_MOUSE,
			Mi: w32.MOUSEINPUT{
				Dx:          dx,
				Dy:          dy,
				MouseData:   mouseData,
				DwFlags:     dwFlags,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
	}
	_ = w32.SendInput(input)
}

func getMouseDownFlag(button di.MouseButton) (flag uint32) {
	switch button {
	case di.ButtonLeft:
		flag = w32.MOUSEEVENTF_LEFTDOWN
	case di.ButtonRight:
		flag = w32.MOUSEEVENTF_RIGHTDOWN
	case di.ButtonMiddle:
		flag = w32.MOUSEEVENTF_MIDDLEDOWN
	}
	return
}

func getMouseUpFlag(button di.MouseButton) (flag uint32) {
	switch button {
	case di.ButtonLeft:
		flag = w32.MOUSEEVENTF_LEFTUP
	case di.ButtonRight:
		flag = w32.MOUSEEVENTF_RIGHTUP
	case di.ButtonMiddle:
		flag = w32.MOUSEEVENTF_MIDDLEUP
	}
	return
}

func getMouseWheelFlag(isHorizontal bool) (flag uint32) {
	flag = uint32(w32.MOUSEEVENTF_WHEEL)
	if isHorizontal {
		flag = w32.MOUSEEVENTF_HWHEEL
	}
	return
}

func getKey(scanCode uint32) (txt string) {
	txt, _ = w32.GetKeyNameText(scanCode, false, false)
	txt = strings.ReplaceAll(txt, " ", "")
	return
}
