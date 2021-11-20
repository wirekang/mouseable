package main

import (
	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/hook"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/logic"
	"github.com/wirekang/mouseable/internal/must"
	"github.com/wirekang/mouseable/internal/view"
)

func main() {
	lg.Logf("Start")
	must.Windows()
	defer func() {
		hook.Uninstall()
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

	logic.SetCursorPos = hook.SetCursorPos
	logic.GetCursorPos = hook.GetCursorPos
	logic.AddCursorPos = hook.AddCursorPos
	logic.MouseDown = hook.MouseDown
	logic.MouseUp = hook.MouseUp
	logic.Wheel = hook.Wheel
	hook.OnKey = logic.OnKey
	hook.Install()

	go logic.Loop()

	err := view.Run()
	if err != nil {
		panic(err)
	}
}
