package io

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/lg"
)

const configVersion = "1"

var configDir string
var configFile string

var DI struct {
	SetConfig func(config def.Config)
}

type functionNameKeyMap map[string]def.FunctionKey
type dataNameValueMap map[string]def.DataValue

type jsonHolder struct {
	Function functionNameKeyMap
	Data     dataNameValueMap
}

func Init() {
	configDir = os.Getenv("APPDATA") + "\\mouseable"
	if cnst.IsDev {
		configDir += "_dev"
	}
	_ = os.Mkdir(configDir, os.ModeDir)
	configFile = configDir + "\\config_v" + configVersion + ".json"
	lg.Logf("ConfigFile: %s", configFile)
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	DI.SetConfig(config)
	return
}

func SaveConfig(config def.Config) (err error) {
	err = saveConfig(config)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	DI.SetConfig(config)
	return
}

func SaveConfigJSON(json string) (err error) {
	c, err := loadConfig([]byte(json))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = saveConfig(c)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func saveConfig(config def.Config) (err error) {
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
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	config, err = loadConfig(bytes)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func loadConfig(bytes []byte) (config def.Config, err error) {
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
