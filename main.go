package main

import (
	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/must"
	"github.com/wirekang/mouseable/internal/view"
)

func main() {
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

	err := view.Init()
	if err != nil {
		panic(err)
	}
}
