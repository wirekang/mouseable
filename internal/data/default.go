package data

import (
	"github.com/wirekang/mouseable/internal/def"
)

func makeDefaultConfig() def.Config {
	return def.Config{
		HotKeyMap: map[*def.HotKeyDef]def.HotKey{
			def.Activate: {
				IsAlt:   true,
				KeyCode: 74,
			},
		},
		FunctionKeyCodeMap: map[*def.FunctionDef]uint32{
			def.Deactivate:  71,
			def.MoveRight:   76,
			def.MoveUp:      75,
			def.MoveLeft:    72,
			def.MoveDown:    74,
			def.ClickLeft:   65,
			def.ClickRight:  68,
			def.ClickMiddle: 83,
			def.WheelUp:     85,
			def.WheelDown:   78,
			def.SniperMode:  32,
			def.Flash:       70,
		},
		DataValueMap: map[*def.DataDef]float64{
			def.Acceleration:    4.0,
			def.Friction:        3.6,
			def.SniperModeSpeed: 1.0,
			def.WheelAmount:     40,
			def.FlashDistance:   300,
		},
	}
}
