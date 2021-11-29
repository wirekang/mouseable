package def

var DefaultConfig = Config{
	FunctionMap: FunctionMap{
		Activate: {
			IsAlt:   true,
			KeyCode: 74,
		},
		Deactivate:      {KeyCode: 71},
		MoveRight:       {KeyCode: 76},
		MoveUp:          {KeyCode: 75},
		MoveLeft:        {KeyCode: 72},
		MoveDown:        {KeyCode: 74},
		ClickLeft:       {KeyCode: 65},
		ClickRight:      {KeyCode: 68},
		ClickMiddle:     {KeyCode: 83},
		WheelUp:         {KeyCode: 85},
		WheelDown:       {KeyCode: 78},
		SniperMode:      {KeyCode: 32},
		TeleportForward: {KeyCode: 70},
	},
	DataMap: DataMap{
		CursorAcceleration: "4.0",
		CursorFriction:     "3.6",
		WheelAcceleration:  "40",
		WheelFriction:      "30",
		SniperModeSpeed:    "1",
		TeleportDistance:   "300",
	},
}
