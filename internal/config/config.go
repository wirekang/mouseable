package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func Load() (cfg Config, err error) {
	fmt.Println("load: " + ConfigPath)
	b, err := os.ReadFile(ConfigPath)
	isNew := false
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("no config file at " + ConfigPath)
			err = nil
			isNew = true
		} else {
			err = errors.Wrap(err, "open config file")
			return
		}

	}

	var config Config
	if isNew {
		fmt.Println("write default config file at " + ConfigPath)
		config = DefaultConfig
		var b []byte
		b, err = jsonMarshal(config)
		if err != nil {
			err = errors.Wrap(err, "marshal default config")
			return
		}

		err = os.WriteFile(ConfigPath, b, 0755)
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

func (c Config) ToJSON() ([]byte, error) {
	return json.Marshal(c)
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
