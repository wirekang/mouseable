package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func Load(path string) (cfg Config, err error) {
	path = strings.ReplaceAll(path, "\\", "/")
	fmt.Println("load: " + path)
	bytes, err := os.ReadFile(path)
	isNew := false
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("no config file at " + path)
			err = nil
			isNew = true
		} else {
			err = errors.Wrap(err, "open config file")
			return
		}

	}

	var config Config
	if isNew {
		fmt.Println("write default config file at " + path)
		config = DefaultConfig
		var bytes []byte
		bytes, err = jsonMarshal(config)
		if err != nil {
			err = errors.Wrap(err, "marshal default config")
			return
		}

		err = os.WriteFile(path, bytes, 0755)
		if err != nil {
			err = errors.Wrap(err, "write default config")
			return
		}

	} else {
		err = json.Unmarshal(bytes, &config)
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
