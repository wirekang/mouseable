package data

import (
	"github.com/wirekang/mouseable/internal/def"
)

func makeDefaultConfig() def.Config {
	return def.Config{
		FunctionKeyCodeMap: map[*def.Function]uint32{
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
		},
		DataValueMap: map[*def.Data]float64{
			def.Acceleration:    5.2,
			def.Friction:        4.0,
			def.SniperModeSpeed: 2.4,
			def.WheelAmount:     40,
		},
		ActivateKey: def.HotKey{
			IsAlt:   true,
			KeyCode: 74,
		},
		DeactivateKey: def.HotKey{
			KeyCode: 186,
		},
	}
}
