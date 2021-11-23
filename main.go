package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/data"
	_ "github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/logic"
	"github.com/wirekang/mouseable/internal/must"
	"github.com/wirekang/mouseable/internal/view"
)

var VERSION string

func main() {
	// checking -dev.exe instead of -dev is due to bug of air.
	// https://github.com/cosmtrek/air/issues/207
	if len(os.Args) == 2 && (os.Args[1] == "-dev.exe" || os.Args[1] == "-dev") {
		lg.IsDev = true
	}

	lg.Logf("Start")
	must.Windows()
	defer func() {
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
	data.Init()
	go logic.Loop()
	view.Run()
}
