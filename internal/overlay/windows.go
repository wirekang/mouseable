//go:build windows

package overlay

import (
	"syscall"
	"unsafe"

	"github.com/JamesHovious/w32"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func New() typ.OverlayManager {
	return &manager{
		cursorWidth:  w32.GetSystemMetrics(w32.SM_CXCURSOR),
		cursorHeight: w32.GetSystemMetrics(w32.SM_CYCURSOR),
		hwnd:         newWindow(),
	}
}

type manager struct {
	isVisible    bool
	cursorWidth  int
	cursorHeight int
	hwnd         w32.HWND
}

func (m *manager) SetPosition(x, y int) {
	if m.isVisible {
		w32.MoveWindow(m.hwnd, x+m.cursorWidth, y+m.cursorHeight, 8, 8, true)
	}
}

func (m *manager) SetVisibility(b bool) {
	m.isVisible = b
}

func (m *manager) Show() {
	if m.isVisible {
		w32.ShowWindow(m.hwnd, w32.SW_SHOWNORMAL)
	}
}

func (m *manager) Hide() {
	w32.ShowWindow(m.hwnd, w32.SW_HIDE)
}

func newWindow() w32.HWND {
	className, _ := syscall.UTF16PtrFromString("MouseableOverlayClassName")
	windowName, _ := syscall.UTF16PtrFromString("MouseableOverlayWindowName")

	var class w32.WNDCLASSEX
	class.Background = newSolidBrush()
	class.ClassName = className
	class.Size = uint32(unsafe.Sizeof(class))
	class.WndProc = syscall.NewCallback(cb)
	w32.RegisterClassEx(&class)
	return w32.CreateWindowEx(
		w32.WS_EX_TOOLWINDOW|w32.WS_EX_TOPMOST|w32.WS_EX_NOACTIVATE,
		className,
		windowName,
		w32.WS_POPUP,
		0, 0, 1, 1,
		0, 0, 0, nil,
	)
}

func newSolidBrush() w32.HBRUSH {
	dll := syscall.NewLazyDLL("Gdi32.dll")
	f := dll.NewProc("CreateSolidBrush")
	b, _, err := f.Call(0x000000FF)
	if err != nil {
		lg.Errorf("CreateSolidBrush: %v", err)
	}
	return w32.HBRUSH(b)
}

func cb(
	hwnd w32.HWND, msg uint32, wparam w32.WPARAM, lparam w32.LPARAM,
) w32.LRESULT {
	return w32.LRESULT(
		w32.DefWindowProc(
			hwnd, msg, uintptr(wparam),
			uintptr(lparam),
		),
	)
}
