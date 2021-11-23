package logic

import (
	"github.com/wirekang/mouseable/internal/def"
)

func SetConfig(config def.Config) {
	mutex.Lock()
	defer mutex.Unlock()

	dataMap = config.DataValueMap
	keyCodeLogicMap = make(
		map[uint32]*logicDef, len(config.FunctionKeyCodeMap),
	)
	for fnc, keyCode := range config.FunctionKeyCodeMap {
		for _, lgc := range logicDefs {
			if lgc.function == fnc {
				keyCodeLogicMap[keyCode] = lgc
			}
		}
	}
}
