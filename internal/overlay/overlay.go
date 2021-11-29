package overlay

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/JamesHovious/w32"

	"github.com/wirekang/mouseable/internal/lg"
)

var hwnd w32.HWND
var state struct {
	sync.Mutex
	isActivating bool
}

func OnCursorMove() {
	hideWindow()
}

func OnCursorStop() {
	state.Lock()
	if state.isActivating {
		showWindow()
	}
	state.Unlock()
}

func OnActivated() {
	state.Lock()
	state.isActivating = true
	state.Unlock()
	showWindow()
}

func OnDeactivated() {
	state.Lock()
	state.isActivating = false
	state.Unlock()
	hideWindow()
}

func Init() {
	initWindow()
}

func initWindow() {
	className, _ := syscall.UTF16PtrFromString("TestClassName")
	windowName, _ := syscall.UTF16PtrFromString("TestWindowName")

	var class w32.WNDCLASSEX
	class.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	class.Background = createSolidBrush()
	class.ClassName = className
	class.Size = uint32(unsafe.Sizeof(class))
	class.WndProc = syscall.NewCallback(cb)
	w32.RegisterClassEx(&class)
	hwnd = w32.CreateWindowEx(
		w32.WS_EX_TOOLWINDOW|w32.WS_EX_TOPMOST|w32.WS_EX_NOACTIVATE,
		className,
		windowName,
		w32.WS_POPUP,
		500, 500, 16, 16,
		0, 0, 0, nil,
	)
	lg.Logf("Overlay HWND %d\n", hwnd)
}

func showWindow() {
	lg.Logf("Show Overlay")
	cursorWidth := w32.GetSystemMetrics(w32.SM_CXCURSOR)
	cursorHeight := w32.GetSystemMetrics(w32.SM_CYCURSOR)
	cursorX, cursorY, _ := w32.GetCursorPos()
	w32.SetWindowPos(
		hwnd, 0, cursorX+cursorWidth-8, cursorY+cursorHeight-8, 16, 16,
		w32.SWP_NOSIZE|w32.SWP_NOACTIVATE,
	)
	w32.ShowWindow(hwnd, w32.SW_SHOWNORMAL)
}

func hideWindow() {
	lg.Logf("Hide Overlay")
	w32.ShowWindow(hwnd, w32.SW_HIDE)
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

func createSolidBrush() w32.HBRUSH {
	dll := syscall.NewLazyDLL("Gdi32.dll")
	f := dll.NewProc("CreateSolidBrush")
	b, _, err := f.Call(0x000000FF)
	if err != nil {
		lg.Errorf("CreateSolidBrush: %v", err)
	}
	return w32.HBRUSH(b)
}
