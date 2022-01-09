package io

import (
	"github.com/wirekang/mouseable/internal/def"
)

func functionMapToNameMap(m def.FunctionMap) (rst functionNameKeyMap) {
	rst = make(functionNameKeyMap, len(m))
	for fnc := range m {
		rst[fnc.Name] = m[fnc]
	}
	return
}

func dataMapToNameMap(m def.DataMap) (rst dataNameValueMap) {
	rst = make(dataNameValueMap, len(m))
	for data := range m {
		rst[data.Name] = m[data]
	}
	return
}

func functionMapFromNameMap(m functionNameKeyMap) (rst def.FunctionMap) {
	rst = make(def.FunctionMap, len(def.FunctionDefinitions))
	for name, key := range m {
		rst[def.FunctionNameMap[name]] = key
	}
	for i := range def.FunctionDefinitions {
		_, ok := rst[def.FunctionDefinitions[i]]
		if !ok {
			rst[def.FunctionDefinitions[i]] = def.FunctionKey{}
		}
	}
	return
}

func dataMapFromNameMap(m dataNameValueMap) (rst map[*def.DataDefinition]def.DataValue) {
	rst = def.DefaultConfig.DataMap
	for name, value := range m {
		rst[def.DataNameMap[name]] = value
	}
	return
}
