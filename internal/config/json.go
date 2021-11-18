package config

import (
	"bytes"
	"encoding/json"

	"github.com/wirekang/mouseable/internal/key"
)

type ConfigJson struct {
	Speed SpeedJson `json:"speed"`
	Key   KeyJson   `json:"key"`
}

type SpeedJson struct {
	Default int `json:"default"`
	Speed1  int `json:"speed_1"`
	Speed2  int `json:"speed_2"`
	Speed3  int `json:"speed_3"`
}

type KeyJson struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
	Right      string `json:"right"`
	RightUp    string `json:"right_up"`
	Up         string `json:"up"`
	LeftUp     string `json:"leftUp"`
	Left       string `json:"left"`
	LeftDown   string `json:"left_down"`
	Down       string `json:"down"`
	RightDown  string `json:"right_down"`
	Speed1     string `json:"speed_1"`
	Speed2     string `json:"speed_2"`
	Speed3     string `json:"speed_3"`
}

func (c ConfigJson) Marshal() ([]byte, error) {
	return jsonMarshal(c)
}

func (c ConfigJson) Config() Config {
	return Config{
		Speed: Speed(c.Speed),
		Key: Key{
			Activate:   key.Parse(c.Key.Activate),
			Deactivate: key.Parse(c.Key.Deactivate),
			Right:      key.Parse(c.Key.Right),
			RightUp:    key.Parse(c.Key.RightUp),
			Up:         key.Parse(c.Key.Up),
			LeftUp:     key.Parse(c.Key.LeftUp),
			Left:       key.Parse(c.Key.Left),
			LeftDown:   key.Parse(c.Key.LeftDown),
			Down:       key.Parse(c.Key.Down),
			RightDown:  key.Parse(c.Key.RightDown),
			Speed1:     key.Parse(c.Key.Speed1),
			Speed2:     key.Parse(c.Key.Speed2),
			Speed3:     key.Parse(c.Key.Speed3),
		},
	}
}

func parse(c Config) (j ConfigJson) {
	j.Speed = SpeedJson(c.Speed)
	j.Key = KeyJson{
		Activate:   c.Key.Activate.String(),
		Deactivate: c.Key.Deactivate.String(),
		Right:      c.Key.Right.String(),
		RightUp:    c.Key.RightUp.String(),
		Up:         c.Key.Up.String(),
		LeftUp:     c.Key.LeftUp.String(),
		Left:       c.Key.Left.String(),
		LeftDown:   c.Key.LeftDown.String(),
		Down:       c.Key.Down.String(),
		RightDown:  c.Key.RightDown.String(),
		Speed1:     c.Key.Speed1.String(),
		Speed2:     c.Key.Speed2.String(),
		Speed3:     c.Key.Speed3.String(),
	}
	return
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
