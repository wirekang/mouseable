package config

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/winsvc/internal/lg"
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

	var config Config
	if isNew {
		lg.Logf("write default config file: %s ", FilePath)
		config = DefaultConfig
		var b []byte
		b, err = config.JSON()
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
		err = json.Unmarshal(b, &config)
		if err != nil {
			err = errors.Wrap(err, "json unmarshal")
			return
		}

	}

	return
}

func (c Config) JSON() ([]byte, error) {
	return jsonMarshal(c)
}

func jsonMarshal(v interface{}) (rst []byte, err error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(v)
	rst = buffer.Bytes()
	return
}
