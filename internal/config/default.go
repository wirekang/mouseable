package config

import (
	"os"
	"path"
)

var DefaultConfigPath = path.Join(os.Getenv("USERPROFILE"), "mouseable.json")

var DefaultConfig = Config{
	Speed: Speed{
		Default: 30,
		Speed1:  5,
		Speed2:  0,
		Speed3:  0,
	},
	Shortcut: Shortcut{
		Activate:   "<A-0>",
		Deactivate: "<A-0>",
		Right:      "H",
		RightUp:    "",
		Up:         "K",
		LeftUp:     "",
		Left:       "L",
		LeftDown:   "",
		Down:       "J",
		RightDown:  "",
		Speed1:     "<Shift>",
		Speed2:     "",
		Speed3:     "",
	},
}
