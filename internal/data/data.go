package data

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

var DI struct {
	SetData func(map[string][]uint32, map[string]string)
}

type jsonHolder struct {
	Keymap map[string][]uint32
	Data   map[string]string
}

func Init() (err error) {
	k, d, err := LoadData()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	DI.SetData(k, d)
	return
}

func SaveData(
	keymap map[string][]uint32, data map[string]string,
) (err error) {
	DI.SetData(keymap, data)

	_ = os.MkdirAll(configDir, os.ModeDir)
	jh := jsonHolder{
		Keymap: keymap,
		Data:   data,
	}
	bytes, err := json.Marshal(jh)
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

func LoadData() (
	keymap map[string][]uint32, data map[string]string, err error,
) {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			keymap = makeDefaultKeymap()
			data = makeDefaultData()
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	var jh jsonHolder
	err = json.Unmarshal(bytes, &jh)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	keymap = jh.Keymap
	data = jh.Data
	return
}

var configDir = os.Getenv("APPDATA") + "\\mouseable"
var configFile = configDir + "\\config.json"
