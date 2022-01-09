package overlay

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/JamesHovious/w32"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/lg"
)

func SetConfig(config def.Config) {
	state.Lock()
	state.isEnabled = config.DataMap[def.ShowOverlay].Bool()
	state.Unlock()
}

var hwnd w32.HWND
var state struct {
	sync.Mutex
	isActivating bool
	isEnabled    bool
	cursorWidth  int
	cursorHeight int
}

func OnCursorMove(x, y int) {
	state.Lock()
	w32.SetWindowPos(hwnd, 0, x+state.cursorWidth-8, y+state.cursorHeight-8, 8, 8, 0)
	state.Unlock()
}

func OnActivated() {
	state.Lock()
	state.isActivating = true
	state.cursorWidth = w32.GetSystemMetrics(w32.SM_CXCURSOR)
	state.cursorHeight = w32.GetSystemMetrics(w32.SM_CYCURSOR)
	if state.isEnabled {
		showWindow()
	}
	state.Unlock()
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
}

func showWindow() {
	w32.ShowWindow(hwnd, w32.SW_SHOWNORMAL)
}

func hideWindow() {
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
