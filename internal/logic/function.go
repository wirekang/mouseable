package logic

type function struct {
	name         string
	keyCodes     []uint32
	onStep       func()
	onActivate   func()
	onDeactivate func()
	isActivated  bool
}

var functions = []*function{
	{
		name:     "MoveRight",
		keyCodes: []uint32{160, 164, 76},
		onStep: func() {
			xSpeed += speed
		},
	},
	{
		name:     "MoveUp",
		keyCodes: []uint32{160, 164, 75},
		onStep: func() {
			ySpeed -= speed
		},
	},
	{
		name:     "MoveLeft",
		keyCodes: []uint32{160, 164, 72},
		onStep: func() {
			xSpeed -= speed
		},
	},
	{
		name:     "MoveDown",
		keyCodes: []uint32{160, 164, 74},
		onStep: func() {
			ySpeed += speed
		},
	},
	{
		name:     "LeftClick",
		keyCodes: []uint32{160, 164, 65},
		onActivate: func() {
			DI.MouseDown(0)
		},
		onDeactivate: func() {
			DI.MouseUp(0)
		},
	},
	{
		name:     "RightClick",
		keyCodes: []uint32{160, 164, 68},
		onActivate: func() {
			DI.MouseDown(0)
		},
		onDeactivate: func() {
			DI.MouseUp(0)
		},
	},
}
