package logic

import (
	"github.com/wirekang/mouseable/internal/def"
)

type logicState struct {
	fixedSpeed     float64
	speedX, speedY float64
	steppingMap    map[*logicDef]struct{}
}

type logicDef struct {
	function *def.Function
	onStart  func(state *logicState)
	onStep   func(state *logicState)
	onStop   func(state *logicState)
}

var logicDefs = []*logicDef{
	{
		function: def.MoveRight,
		onStep: func(s *logicState) {
			s.speedX += dataMap[def.Acceleration]
		},
	},
	{
		function: def.MoveUp,
		onStep: func(s *logicState) {
			s.speedY -= dataMap[def.Acceleration]
		},
	},
	{
		function: def.MoveLeft,
		onStep: func(s *logicState) {
			s.speedX -= dataMap[def.Acceleration]
		},
	},
	{
		function: def.MoveDown,
		onStep: func(s *logicState) {
			s.speedY += dataMap[def.Acceleration]
		},
	},
	{
		function: def.ClickLeft,
		onStart: func(_ *logicState) {
			DI.MouseDown(0)
		},
		onStop: func(_ *logicState) {
			DI.MouseUp(0)
		},
	},
	{
		function: def.ClickRight,
		onStart: func(_ *logicState) {
			DI.MouseDown(1)
		},
		onStop: func(_ *logicState) {
			DI.MouseUp(1)
		},
	},
	{
		function: def.ClickMiddle,
		onStart: func(_ *logicState) {
			DI.MouseDown(2)
		},
		onStop: func(_ *logicState) {
			DI.MouseUp(2)
		},
	},
	{
		function: def.WheelUp,
		onStep: func(_ *logicState) {
			DI.Wheel(int(dataMap[def.WheelAmount]), false)
		},
	},
	{
		function: def.WheelDown,
		onStep: func(_ *logicState) {
			DI.Wheel(-int(dataMap[def.WheelAmount]), false)
		},
	},
	{
		function: def.SniperMode,
		onStart: func(s *logicState) {
			s.fixedSpeed = dataMap[def.SniperModeSpeed]
		},
		onStop: func(s *logicState) {
			s.fixedSpeed = 0
		},
	},
}
