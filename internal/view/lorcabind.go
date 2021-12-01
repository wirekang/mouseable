package view

import (
	"os/exec"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/def"
)

func __loadBind__() interface{} {
	configHolder.Lock()
	defer configHolder.Unlock()
	m := make(map[string]interface{})
	m["functionDefinitions"] = def.FunctionDefinitions
	m["dataDefinitions"] = def.DataDefinitions
	fnm := make(map[string]def.FunctionKey, len(configHolder.FunctionMap))
	for fd := range configHolder.FunctionMap {
		fnm[fd.Name] = configHolder.FunctionMap[fd]
	}
	m["functionNameKeyMap"] = fnm
	dnm := make(map[string]def.DataValue, len(configHolder.DataMap))
	for dd := range configHolder.DataMap {
		dnm[dd.Name] = configHolder.DataMap[dd]
	}
	m["dataNameValueMap"] = dnm
	m["version"] = cnst.VERSION
	return m
}

func __getKeyCode__() uint32 {
	DI.NormalKeyChan <- 0
	return <-DI.NormalKeyChan
}

func __changeFunction__(name string, key def.FunctionKey) bool {
	configHolder.Lock()
	defer configHolder.Unlock()
	configHolder.FunctionMap[def.FunctionNameMap[name]] = key
	err := DI.SaveConfig(configHolder.Config)
	return err == nil
}

func __changeData__(name string, value string) bool {
	configHolder.Lock()
	defer configHolder.Unlock()
	configHolder.DataMap[def.DataNameMap[name]] = def.DataValue(value)
	err := DI.SaveConfig(configHolder.Config)
	return err == nil
}

func __openGitHub__() {
	exec.Command("rundll32", "url.dll,FileProtocolHandler", "https://github.com/wirekang/mouseable").Start()
}
