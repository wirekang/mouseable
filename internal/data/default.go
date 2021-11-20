package data

func makeDefaultKeymap() map[string][]uint32 {
	return map[string][]uint32{
		"Activate":    {164, 76},
		"MoveRight":   {164, 76},
		"MoveUp":      {164, 75},
		"MoveLeft":    {164, 72},
		"MoveDown":    {164, 74},
		"LeftClick":   {164, 65},
		"RightClick":  {164, 68},
		"MiddleClick": {164, 83},
		"WheelUp":     {164, 85},
		"WheelDown":   {164, 78},
	}
}

func makeDefaultData() map[string]string {
	return map[string]string{
		"acceleration": "4",
		"friction":     "3",
		"wheelAmount":  "10",
	}
}
