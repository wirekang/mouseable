package data

func makeDefaultKeymap() map[string][]uint32 {
	return map[string][]uint32{
		"Activate":    {164, 74},
		"Deactivate":  {186},
		"MoveRight":   {76},
		"MoveUp":      {75},
		"MoveLeft":    {72},
		"MoveDown":    {74},
		"LeftClick":   {65},
		"RightClick":  {68},
		"MiddleClick": {83},
		"WheelUp":     {85},
		"WheelDown":   {78},
		"Sniper":      {32},
	}
}

func makeDefaultData() map[string]string {
	return map[string]string{
		"acceleration": "5.0",
		"friction":     "4.0",
		"sniper":       "3.0",
		"wheelAmount":  "30",
	}
}
