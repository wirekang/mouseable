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
		name: "MoveRight",
		onStep: func() {
			xSpeed += getInt("acceleration")
		},
	},
	{
		name: "MoveUp",
		onStep: func() {
			ySpeed -= getInt("acceleration")
		},
	},
	{
		name: "MoveLeft",
		onStep: func() {
			xSpeed -= getInt("acceleration")
		},
	},
	{
		name: "MoveDown",
		onStep: func() {
			ySpeed += getInt("acceleration")
		},
	},
	{
		name: "LeftClick",
		onActivate: func() {
			DI.MouseDown(0)
		},
		onDeactivate: func() {
			DI.MouseUp(0)
		},
	},
	{
		name: "RightClick",
		onActivate: func() {
			DI.MouseDown(1)
		},
		onDeactivate: func() {
			DI.MouseUp(1)
		},
	},
	{
		name: "MiddleClick",
		onActivate: func() {
			DI.MouseDown(2)
		},
		onDeactivate: func() {
			DI.MouseUp(2)
		},
	},
	{
		name: "WheelUp",
		onStep: func() {
			DI.Wheel(getInt("wheelAmount"), false)
		},
	},
	{
		name: "WheelDown",
		onStep: func() {
			DI.Wheel(-getInt("wheelAmount"), false)
		},
	},
}
