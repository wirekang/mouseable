package data

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/def"
)

var DI struct {
	SetConfig func(config def.Config)
}

type jsonHolder struct {
	HotKeyNameMap   map[string]def.HotKey
	FunctionNameMap map[string]uint32
	DataNameMap     map[string]float64
}

func Init() {
	config := LoadConfig()
	DI.SetConfig(config)
	return
}

func SaveConfig(config def.Config) {
	saveData(config)
	DI.SetConfig(config)
}

func saveData(config def.Config) (err error) {
	DI.SetConfig(config)
	_ = os.MkdirAll(configDir, os.ModeDir)
	jh := jsonHolder{
		HotKeyNameMap:   hotKeyMapToNameMap(config.HotKeyMap),
		FunctionNameMap: functionMapToNameMap(config.FunctionKeyCodeMap),
		DataNameMap:     dataMapToNameMap(config.DataValueMap),
	}
	bytes, err := json.MarshalIndent(jh, "", "    ")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = os.WriteFile(configFile, bytes, 0755)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func LoadConfig() (config def.Config) {
	config, err := loadConfig()
	if err != nil {
		config = makeDefaultConfig()
	}
	return
}

func loadConfig() (config def.Config, err error) {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	var jh jsonHolder
	err = json.Unmarshal(bytes, &jh)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	config = def.Config{
		HotKeyMap:          nameMapToHotKeyMap(jh.HotKeyNameMap),
		FunctionKeyCodeMap: nameMapToFunctionMap(jh.FunctionNameMap),
		DataValueMap:       nameMapToDataMap(jh.DataNameMap),
	}

	return
}

func hotKeyMapToNameMap(m map[*def.HotKeyDef]def.HotKey) (rst map[string]def.HotKey) {
	rst = make(map[string]def.HotKey, len(m))
	for hk := range m {
		rst[hk.Name] = m[hk]
	}
	return
}

func functionMapToNameMap(m map[*def.FunctionDef]uint32) (rst map[string]uint32) {
	rst = make(map[string]uint32, len(m))
	for fnc := range m {
		rst[fnc.Name] = m[fnc]
	}
	return
}

func dataMapToNameMap(m map[*def.DataDef]float64) (rst map[string]float64) {
	rst = make(map[string]float64, len(m))
	for data := range m {
		rst[data.Name] = m[data]
	}
	return
}

func nameMapToHotKeyMap(m map[string]def.HotKey) (rst map[*def.HotKeyDef]def.HotKey) {
	rst = makeDefaultConfig().HotKeyMap
	for name, keyCode := range m {
		rst[def.HotKeyNameMap[name]] = keyCode
	}
	return
}

func nameMapToFunctionMap(m map[string]uint32) (rst map[*def.FunctionDef]uint32) {
	rst = makeDefaultConfig().FunctionKeyCodeMap
	for name, keyCode := range m {
		rst[def.FunctionNameMap[name]] = keyCode
	}
	return
}

func nameMapToDataMap(m map[string]float64) (rst map[*def.DataDef]float64) {
	rst = makeDefaultConfig().DataValueMap
	for name, value := range m {
		rst[def.DataNameMap[name]] = value
	}
	return
}

var configDir = os.Getenv("APPDATA") + "\\mouseable"
var configFile = configDir + "\\config.json"
