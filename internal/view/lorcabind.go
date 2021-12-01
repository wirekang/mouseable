package view

import (
	"github.com/wirekang/mouseable/internal/def"
)

func __loadBind__() interface{} {
	m := make(map[string]interface{})
	m["functionDefinitions"] = def.FunctionDefinitions
	m["dataDefinitions"] = def.DataDefinitions
	fnm := make(map[string]def.FunctionKey, len(config.FunctionMap))
	for fd := range config.FunctionMap {
		fnm[fd.Name] = config.FunctionMap[fd]
	}
	m["functionNameKeyMap"] = fnm
	dnm := make(map[string]def.DataValue, len(config.DataMap))
	for dd := range config.DataMap {
		dnm[dd.Name] = config.DataMap[dd]
	}
	m["dataNameValueMap"] = dnm
	return m
}

func __getKeyCode__() uint32 {
	DI.NormalKeyChan <- 0
	return <-DI.NormalKeyChan
}

func __changeFunction__(name string, key def.FunctionKey) bool {
	return true
}
