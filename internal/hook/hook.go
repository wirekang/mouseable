package hook

import (
	"unsafe"

	"github.com/JamesHovious/w32"
)

var hHook w32.HHOOK
var DI struct {
	OnKey func(keyCode uint32, isDown bool) (preventDefault bool)
}

func Install() {
	hHook = w32.SetWindowsHookEx(w32.WH_KEYBOARD_LL, hookProc, 0, 0)
}

func Uninstall() {
	w32.UnhookWindowsHookEx(hHook)
}

func SetCursorPos(x, y int) {
	w32.SetCursorPos(x, y)
}

func AddCursorPos(dx, dy int32) {
	sendMouseInput(dx, dy, 0, w32.MOUSEEVENTF_MOVE)
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

var hookProc w32.HOOKPROC = func(
	code int, wParam w32.WPARAM, lParam w32.LPARAM,
) w32.LRESULT {
	data := *(*w32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	flagMap := map[string]w32.DWORD{
		"Extended":      data.Flags & (w32.KF_EXTENDED >> 8),
		"LowerInjected": data.Flags & 0x00000002,
		"Injected":      data.Flags & 0x00000010,
		"AltDown":       data.Flags & (w32.KF_ALTDOWN >> 8),
		"Up":            data.Flags & (w32.KF_UP >> 8),
	}

	if DI.OnKey != nil {
		isDown := flagMap["Up"] == 0
		preventDefault := DI.OnKey(uint32(data.VkCode), isDown)

		if preventDefault {
			return 1
		}
	}

	return w32.CallNextHookEx(0, code, wParam, lParam)
}
