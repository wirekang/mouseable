package view

import (
	"sync"

	"github.com/wirekang/mouseable/internal/def"
)

var DI struct {
	LoadConfig     func() (def.Config, error)
	SaveConfigJSON func(string) error
	GetKeyText     func(keyCode uint32) (string, bool)
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
	go openUI()
	runNotifyIcon()
}
