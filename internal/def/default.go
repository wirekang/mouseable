package def

var DefaultConfig = Config{
	FunctionMap: FunctionMap{
		Activate: {
			IsAlt:   true,
			KeyCode: 74,
		},
		Deactivate:      {KeyCode: 186},
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
		DoublePressSpeed:    "300",
		CursorAccelerationH: "2.8",
		CursorAccelerationV: "2.8",
		CursorFrictionH:     "2.5",
		CursorFrictionV:     "2.5",
		WheelAccelerationH:  "5",
		WheelAccelerationV:  "5",
		WheelFrictionH:      "4",
		WheelFrictionV:      "4",
		SniperModeSpeedH:    "1",
		SniperModeSpeedV:    "1",
		TeleportDistanceF:   "300",
		TeleportDistanceH:   "300",
		TeleportDistanceV:   "300",
	},
}
