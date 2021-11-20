package logic

type function struct {
	name        string
	keyCodes    []uint32
	step        func()
	isActivated bool
}

var functions = []*function{
	{
		name:     "MoveRight",
		keyCodes: []uint32{164, 76},
		step: func() {
			xSpeed += speed
		},
	},
	{
		name:     "MoveUp",
		keyCodes: []uint32{164, 75},
		step: func() {
			ySpeed -= speed
		},
	},
	{
		name:     "MoveLeft",
		keyCodes: []uint32{164, 72},
		step: func() {
			xSpeed -= speed
		},
	},
	{
		name:     "MoveDown",
		keyCodes: []uint32{164, 74},
		step: func() {
			ySpeed += speed
		},
	},
}
