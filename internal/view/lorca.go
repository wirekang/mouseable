package view

import (
	"os"

	"github.com/pkg/errors"
	"github.com/zserge/lorca"
)

func mustChrome() {
	if lorca.LocateChrome() == "" {
		AlertError("Google Chrome not found. Mouseable can't render GUI.")
		os.Exit(1)
	}
}

func waitLorca() (err error) {
	ui, err := lorca.New("http://"+host, "", 800, 800)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer ui.Close()
	<-ui.Done()
	return
}
