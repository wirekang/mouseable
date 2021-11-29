package io

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/def"
)

const configVersion = "1"

var configDir = os.Getenv("APPDATA") + "\\mouseable"
var configFile = configDir + "\\config_v" + configVersion + ".json"

var DI struct {
	SetConfig func(config def.Config)
}

type functionNameMap map[string]def.FunctionKey
type dataNameMap map[string]def.DataValue

type jsonHolder struct {
	Function functionNameMap
	Data     dataNameMap
}

func Init() {
	config := LoadConfig()
	DI.SetConfig(config)
	return
}

func SaveConfig(config def.Config) {
	err := saveData(config)
	if err != nil {
		panic(err)
	}

	DI.SetConfig(config)
}

func saveData(config def.Config) (err error) {
	DI.SetConfig(config)
	_ = os.MkdirAll(configDir, os.ModeDir)
	jh := jsonHolder{
		Function: functionMapToNameMap(config.FunctionMap),
		Data:     dataMapToNameMap(config.DataMap),
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
		config = def.DefaultConfig
	}
	return
}

func loadConfig() (config def.Config, err error) {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	var nameConfig jsonHolder
	err = json.Unmarshal(bytes, &nameConfig)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	config = def.Config{
		FunctionMap: functionMapFromNameMap(nameConfig.Function),
		DataMap:     dataMapFromNameMap(nameConfig.Data),
	}

	return
}

func functionMapToNameMap(m def.FunctionMap) (rst functionNameMap) {
	rst = make(functionNameMap, len(m))
	for fnc := range m {
		rst[fnc.Name] = m[fnc]
	}
	return
}

func dataMapToNameMap(m def.DataMap) (rst dataNameMap) {
	rst = make(dataNameMap, len(m))
	for data := range m {
		rst[data.Name] = m[data]
	}
	return
}

func functionMapFromNameMap(m functionNameMap) (rst def.FunctionMap) {
	rst = def.DefaultConfig.FunctionMap
	for name, key := range m {
		rst[def.GetFunctionDefinitionByName(name)] = key
	}
	return
}

func dataMapFromNameMap(m dataNameMap) (rst map[*def.DataDefinition]def.DataValue) {
	rst = def.DefaultConfig.DataMap
	for name, value := range m {
		rst[def.GetDataDefinitionByName(name)] = value
	}
	return
}
