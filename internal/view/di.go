package view

import (
	"github.com/wirekang/mouseable/internal/def"
)

var DI struct {
	LoadConfig func() def.Config
	SaveConfig func(def.Config)
}
