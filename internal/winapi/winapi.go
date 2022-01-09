package winapi

import (
	"fmt"
	"unsafe"

	"github.com/JamesHovious/w32"
)

var DI struct {
	OnKey        func(keyCode uint32, isDown bool) (preventDefault bool)
	OnCursorMove func(x, y int)
}

var hHooks []w32.HHOOK

func Hook() {
	f := func(idHook int, lpfn w32.HOOKPROC) {
		hhook := w32.SetWindowsHookEx(idHook, lpfn, 0, 0)
		if hhook == 0 {
			panic(fmt.Sprintf("Hook failed. id: %d, lastErrorCode: %d", idHook, w32.GetLastError()))
		}

		hHooks = append(hHooks, hhook)
	}
	f(w32.WH_KEYBOARD_LL, keyboardProc)
	f(w32.WH_MOUSE_LL, mouseProc)
}

func UnHook() {
	for i := range hHooks {
		w32.UnhookWindowsHookEx(hHooks[i])
	}
}

func SetCursorPos(x, y int) {
	w32.SetCursorPos(x, y)
}

func AddCursorPos(dx, dy int) {
	sendMouseInput(int32(dx), int32(dy), 0, w32.MOUSEEVENTF_MOVE)
}

func GetCursorPos() (x, y int) {
	x, y, _ = w32.GetCursorPos()
	return
}

// MouseDown send mouse down event.
//
// left = 0
//
// right = 1
//
// middle = 2
func MouseDown(button int) {
	var flag uint32
	switch button {
	case 0:
		flag = w32.MOUSEEVENTF_LEFTDOWN
	case 1:
		flag = w32.MOUSEEVENTF_RIGHTDOWN
	case 2:
		flag = w32.MOUSEEVENTF_MIDDLEDOWN
	}
	sendMouseInput(0, 0, 0, flag)
}

// MouseUp send mouse up event. checkout MouseDown for button.
func MouseUp(button int) {
	var flag uint32
	switch button {
	case 0:
		flag = w32.MOUSEEVENTF_LEFTUP
	case 1:
		flag = w32.MOUSEEVENTF_RIGHTUP
	case 2:
		flag = w32.MOUSEEVENTF_MIDDLEUP
	}
	sendMouseInput(0, 0, 0, flag)
}

func Wheel(amount int, horizontal bool) {
	flag := uint32(w32.MOUSEEVENTF_WHEEL)
	if horizontal {
		flag = w32.MOUSEEVENTF_HWHEEL
	}
	sendMouseInput(0, 0, uint32(amount), flag)
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
	w32.SendInput(input)
}

func GetKeyText(keyCode uint32) (txt string, ok bool) {
	scanCode := w32.MapVirtualKeyEx(uint(keyCode), w32.MAPVK_VK_TO_VSC, 0)
	txt, ok = w32.GetKeyNameText(uint32(scanCode), false, false)
	return
}

var keyboardProc w32.HOOKPROC = func(code int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	data := *(*w32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	// fExtended := data.Flags & (w32.KF_EXTENDED >> 8)
	// fLowerInjected := data.Flags & 0x00000002
	// fInjected := data.Flags & 0x00000010
	// fAltDown := data.Flags & (w32.KF_ALTDOWN >> 8)
	fUp := data.Flags & (w32.KF_UP >> 8)

	if DI.OnKey != nil {
		isDown := fUp == 0
		preventDefault := DI.OnKey(uint32(data.VkCode), isDown)

		if preventDefault {
			return 1
		}
	}

	return w32.CallNextHookEx(0, code, wParam, lParam)
}

var mouseProc w32.HOOKPROC = func(code int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	data := *(*w32.MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	if DI.OnCursorMove != nil {
		go DI.OnCursorMove(int(data.Pt.X), int(data.Pt.Y))
	}
	return w32.CallNextHookEx(0, code, wParam, lParam)
}
