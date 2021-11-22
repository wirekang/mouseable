package logic

import (
	"math"
)

type function struct {
	name               string
	keyCodes           []uint32
	onStart            func()
	onStep             func()
	onStop             func()
	isIgnoreDeactivate bool

	// internal
	isStepping bool
}

var functions = []*function{
	{
		name:               "Activate",
		isIgnoreDeactivate: true,
		onStart: func() {
			state.isActivated = true
		},
	},
	{
		name: "Deactivate",
		onStart: func() {
			state.isActivated = false
		},
	},
	{
		name: "HoldToActivate",
		onStart: func() {
			state.isActivated = true
		},
		onStop: func() {
			state.isActivated = false
		},
	},

	{
		name: "MoveRight",
		onStep: func() {
			state.speedX += getFloat("acceleration")
		},
	},
	{
		name: "MoveUp",
		onStep: func() {
			state.speedY -= getFloat("acceleration")
		},
	},
	{
		name: "MoveLeft",
		onStep: func() {
			state.speedX -= getFloat("acceleration")
		},
	},
	{
		name: "MoveDown",
		onStep: func() {
			state.speedY += getFloat("acceleration")
		},
	},
	{
		name: "LeftClick",
		onStart: func() {
			DI.MouseDown(0)
		},
		onStop: func() {
			DI.MouseUp(0)
		},
	},
	{
		name: "RightClick",
		onStart: func() {
			DI.MouseDown(1)
		},
		onStop: func() {
			DI.MouseUp(1)
		},
	},
	{
		name: "MiddleClick",
		onStart: func() {
			DI.MouseDown(2)
		},
		onStop: func() {
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
	{
		name: "Sniper",
		onStep: func() {
			sniper := getFloat("sniper")
			state.speedX = math.Min(sniper, math.Max(state.speedX, -sniper))
			state.speedY = math.Min(sniper, math.Max(state.speedY, -sniper))
		},
	},
}
