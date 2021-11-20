package main

import (
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/data"
	"github.com/wirekang/mouseable/internal/hook"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/logic"
	"github.com/wirekang/mouseable/internal/must"
	"github.com/wirekang/mouseable/internal/view"
)

func main() {
	// checking -dev.exe instead of -dev is due to bug of air.
	// https://github.com/cosmtrek/air/issues/207
	if len(os.Args) == 2 && os.Args[1] == "-dev.exe" {
		lg.IsDev = true
	}

	lg.Logf("Start")
	must.Windows()
	defer func() {
		err := recover()
		if err != nil {
			lg.Errorf("panic: %v", err)
			if st, ok := err.(interface {
				StackTrace() errors.StackTrace
			}); ok {
				lg.Errorf("StackTrace: \n%+v", st.StackTrace())
			}
		}
		lg.Logf("EXIT")
	}()

	logic.DI.SetCursorPos = hook.SetCursorPos
	logic.DI.GetCursorPos = hook.GetCursorPos
	logic.DI.AddCursorPos = hook.AddCursorPos
	logic.DI.MouseDown = hook.MouseDown
	logic.DI.MouseUp = hook.MouseUp
	logic.DI.Wheel = hook.Wheel
	hook.DI.OnKey = logic.OnKey
	view.DI.LoadData = data.LoadData
	view.DI.SaveData = data.SaveData
	data.DI.SetData = logic.SetData

	err := data.Init()
	if err != nil {
		panic(err)
	}

	hook.Install()
	defer hook.Uninstall()

	go logic.Loop()
	err = view.Run()
	if err != nil {
		panic(err)
	}
}
