package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/key"
	"github.com/wirekang/mouseable/internal/lg"
)

func Load() (cfg Config, err error) {
	lg.Logf("load config: %s", FilePath)
	b, err := os.ReadFile(FilePath)
	isNew := false
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			lg.Logf("no config file: %s", FilePath)
			err = nil
			isNew = true
		} else {
			err = errors.Wrap(err, "open config file")
			return
		}
	}

	if isNew {
		lg.Logf("write default config file: %s ", FilePath)
		var b []byte
		b, err = parse(DefaultConfig).Marshal()
		if err != nil {
			err = errors.Wrap(err, "marshal default config")
			return
		}

		err = os.WriteFile(FilePath, b, 0755)
		if err != nil {
			err = errors.Wrap(err, "write default config")
			return
		}

	} else {
		var configJson ConfigJson
		err = json.Unmarshal(b, &configJson)
		if err != nil {
			err = errors.Wrap(err, "json unmarshal")
			return
		}
		cfg = configJson.Config()
	}

	return
}

type Config struct {
	Speed Speed
	Key   Key
}

type Speed struct {
	Default int
	Speed1  int
	Speed2  int
	Speed3  int
}

type Key struct {
	Activate   key.Key
	Deactivate key.Key
	Right      key.Key
	RightUp    key.Key
	Up         key.Key
	LeftUp     key.Key
	Left       key.Key
	LeftDown   key.Key
	Down       key.Key
	RightDown  key.Key
	Speed1     key.Key
	Speed2     key.Key
	Speed3     key.Key
}
