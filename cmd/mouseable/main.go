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
		lg.Logf("EXIT")
	}()

	err := view.Init()
	if err != nil {
		if st, ok := err.(stackTracer); ok {
			st.StackTrace()
		}
		panic(err)
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
