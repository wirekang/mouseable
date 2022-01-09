package view

import (
	"os"

	"github.com/pkg/errors"
	"github.com/zserge/lorca"
)

func bindLorca(ui lorca.UI) (err error) {
	a := func(name string, f interface{}) {
		err := ui.Bind("__"+name+"__", f)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	a(
		"getKeyText",
		func(keyCode uint32) string {
			s, ok := DI.GetKeyText(keyCode)
			if !ok {
				return "<ERROR>"
			}
			return s
		},
	)
	a(
		"terminate",
		func() {
			ui.Close()
			os.Exit(0)
		},
	)
	return
}
