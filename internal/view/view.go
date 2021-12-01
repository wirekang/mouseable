package view

import (
	"sync"

	"github.com/wirekang/mouseable/internal/def"
)

var DI struct {
	LoadConfig    func() (def.Config, error)
	SaveConfig    func(def.Config) error
	NormalKeyChan chan uint32
}

var configHolder struct {
	sync.Mutex
	def.Config
}

func Run() {
	var err error
	configHolder.Config, err = DI.LoadConfig()
	if err != nil {
		panic(err)
	}

	mustChrome()
	go serveAndOpen()
	runNotifyIcon()
}
