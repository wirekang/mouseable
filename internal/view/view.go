package view

import (
	"github.com/wirekang/mouseable/internal/def"
)

var DI struct {
	LoadConfig    func() (def.Config, error)
	SaveConfig    func(def.Config) error
	NormalKeyChan chan uint32
}

func Run() {
	mustChrome()
	go serveAndOpen()
	var err error
	config, err = DI.LoadConfig()
	if err != nil {
		AlertError(err.Error())
	}

	runNotifyIcon()
}
