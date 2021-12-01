package io

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/lg"
)

const configVersion = "1"

var configDir = os.Getenv("APPDATA") + "\\mouseable"
var configFile = configDir + "\\config_v" + configVersion + ".json"

var DI struct {
	SetConfig func(config def.Config)
}

type functionNameKeyMap map[string]def.FunctionKey
type dataNameValueMap map[string]def.DataValue

type jsonHolder struct {
	Function functionNameKeyMap
	Data     dataNameValueMap
}

func Init() (err error) {
	lg.Logf("ConfigFile: %s", configFile)
	config, err := LoadConfig()
	DI.SetConfig(config)
	return
}

func SaveConfig(config def.Config) (err error) {
	err = saveData(config)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	DI.SetConfig(config)
	return
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

func LoadConfig() (config def.Config, err error) {
	config, err = loadConfig()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func loadConfig() (config def.Config, err error) {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
			config = def.Config{
				FunctionMap: functionMapFromNameMap(functionMapToNameMap(def.DefaultConfig.FunctionMap)),
				DataMap:     dataMapFromNameMap(dataMapToNameMap(def.DefaultConfig.DataMap)),
			}
			return
		}

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
