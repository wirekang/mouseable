package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/io"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/logic"
	"github.com/wirekang/mouseable/internal/must"
	"github.com/wirekang/mouseable/internal/overlay"
	"github.com/wirekang/mouseable/internal/view"
	"github.com/wirekang/mouseable/internal/winapi"
)

//go:embed assets
var Asset embed.FS

var VERSION string

func main() {
	lg.Logf("Start")
	cnst.VERSION = VERSION
	cnst.AssetFS = Asset
	// checking -dev.exe instead of -dev is due to bug of air.
	// https://github.com/cosmtrek/air/issues/207
	if len(os.Args) == 2 && (os.Args[1] == "-dev.exe" || os.Args[1] == "-dev") {
		cnst.IsDev = true
	}

	if !io.Lock() {
		view.AlertError("Mouseable is already running.")
		return
	}

	must.Windows()

	defer func() {
		io.Unlock()
		err := recover()
		if err != nil {
			msg := fmt.Sprintf("panic: %v\n\n", err)
			if st, ok := err.(interface {
				StackTrace() errors.StackTrace
			}); ok {
				msg += fmt.Sprintf("StackTrace: \n%+v", st.StackTrace())
			}
			lg.Errorf(msg)
			view.AlertError(msg)
		}
		lg.Logf("EXIT")
	}()

	di.Init()
	io.Init()
	go logic.Loop()
	winapi.Hook()
	defer winapi.UnHook()
	overlay.Init()
	view.Run()
}
