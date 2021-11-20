package logic

import (
	"sync"
)

type function struct {
	name         string
	keyCodes     []uint32
	onStep       func()
	onActivate   func()
	onDeactivate func()
	isActivated  bool
}

var functionsMutex sync.Mutex
var functions = []*function{
	{
		name:     "MoveRight",
		keyCodes: []uint32{164, 76},
		onStep: func() {
			xSpeed += speed
		},
	},
	{
		name:     "MoveUp",
		keyCodes: []uint32{164, 75},
		onStep: func() {
			ySpeed -= speed
		},
	},
	{
		name:     "MoveLeft",
		keyCodes: []uint32{164, 72},
		onStep: func() {
			xSpeed -= speed
		},
	},
	{
		name:     "MoveDown",
		keyCodes: []uint32{164, 74},
		onStep: func() {
			ySpeed += speed
		},
	},
	{
		name:     "LeftClick",
		keyCodes: []uint32{164, 65},
		onActivate: func() {
			DI.MouseDown(0)
		},
		onDeactivate: func() {
			DI.MouseUp(0)
		},
	},
	{
		name:     "RightClick",
		keyCodes: []uint32{164, 68},
		onActivate: func() {
			DI.MouseDown(1)
		},
		onDeactivate: func() {
			DI.MouseUp(1)
		},
	},
	{
		name:     "MiddleClick",
		keyCodes: []uint32{164, 83},
		onActivate: func() {
			DI.MouseDown(2)
		},
		onDeactivate: func() {
			DI.MouseUp(2)
		},
	},
	{
		name:     "WheelUp",
		keyCodes: []uint32{164, 85},
		onStep: func() {
			DI.Wheel(wheelAmount, false)
		},
	},
	{
		name:     "WheelDown",
		keyCodes: []uint32{164, 78},
		onStep: func() {
			DI.Wheel(-wheelAmount, false)
		},
	},
}
