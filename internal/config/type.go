package config

import (
	"github.com/wirekang/winsvc/internal/key"
)

type Config struct {
	Speed    Speed    `json:"speed"`
	Shortcut Shortcut `json:"shortcut"`
}

type Speed struct {
	Default int `json:"default"`
	Speed1  int `json:"speed_1"`
	Speed2  int `json:"speed_2"`
	Speed3  int `json:"speed_3"`
}

type Shortcut struct {
	Activate   key.Key `json:"activate"`
	Deactivate key.Key `json:"deactivate"`
	Right      key.Key `json:"right"`
	RightUp    key.Key `json:"right_up"`
	Up         key.Key `json:"up"`
	LeftUp     key.Key `json:"leftUp"`
	Left       key.Key `json:"left"`
	LeftDown   key.Key `json:"left_down"`
	Down       key.Key `json:"down"`
	RightDown  key.Key `json:"right_down"`
	Speed1     key.Key `json:"speed_1"`
	Speed2     key.Key `json:"speed_2"`
	Speed3     key.Key `json:"speed_3"`
}
