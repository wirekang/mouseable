//go:build windows

package hook

import (
	"fmt"
	"unsafe"

	"github.com/JamesHovious/w32"

	"github.com/wirekang/mouseable/internal/typ"
)

func New() typ.HookManager {
	return &manager{}
}

type manager struct {
	onKeyListener    typ.OnKeyListener
	onCursorListener typ.OnCursorListener
	hhooks           []w32.HHOOK
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

func (m *manager) MouseDown(button typ.MouseButton) {
	sendMouseInput(0, 0, 0, getMouseDownFlag(button))
}

func (m *manager) MouseUp(button typ.MouseButton) {
	sendMouseInput(0, 0, 0, getMouseUpFlag(button))
}

func (m *manager) Wheel(amount int, isHorizontal bool) {
	if amount != 0 {
		sendMouseInput(0, 0, uint32(amount), getMouseWheelFlag(isHorizontal))
	}
}

func (m *manager) SetOnKeyListener(listener typ.OnKeyListener) {
	m.onKeyListener = listener
}

func (m *manager) SetOnCursorListener(listener typ.OnCursorListener) {
	m.onCursorListener = listener
}

func (m *manager) Install() {
	f := func(idHook int, lpfn w32.HOOKPROC) {
		hhook := w32.SetWindowsHookEx(idHook, lpfn, 0, 0)
		if hhook == 0 {
			panic(fmt.Sprintf("hook failed. id: %d", idHook))
			return
		}

		m.hhooks = append(m.hhooks, hhook)
	}
	f(w32.WH_KEYBOARD_LL, m.keyboardProc)
	f(w32.WH_MOUSE_LL, m.mouseProc)
	return
}

func (m *manager) Uninstall() {
	for i := range m.hhooks {
		w32.UnhookWindowsHookEx(m.hhooks[i])
	}
}

func (m *manager) keyboardProc(code int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	data := *(*w32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	// fExtended := data.Flags & (w32.KF_EXTENDED >> 8)
	// fLowerInjected := data.Flags & 0x00000002
	// fInjected := data.Flags & 0x00000010
	// fAltDown := data.Flags & (w32.KF_ALTDOWN >> 8)
	fUp := data.Flags & (w32.KF_UP >> 8)

	isDown := fUp == 0
	txt := getKey(uint32(data.ScanCode))
	if m.onKeyListener(typ.Key(txt), isDown) {
		return 1
	}

	return w32.CallNextHookEx(0, code, wParam, lParam)
}

func (m *manager) mouseProc(code int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	data := *(*w32.MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	m.onCursorListener(int(data.Pt.X), int(data.Pt.Y))
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

func getMouseDownFlag(button typ.MouseButton) (flag uint32) {
	switch button {
	case typ.Left:
		flag = w32.MOUSEEVENTF_LEFTDOWN
	case typ.Right:
		flag = w32.MOUSEEVENTF_RIGHTDOWN
	case typ.Middle:
		flag = w32.MOUSEEVENTF_MIDDLEDOWN
	}
	return
}

func getMouseUpFlag(button typ.MouseButton) (flag uint32) {
	switch button {
	case typ.Left:
		flag = w32.MOUSEEVENTF_LEFTUP
	case typ.Right:
		flag = w32.MOUSEEVENTF_RIGHTUP
	case typ.Middle:
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
	return
}
