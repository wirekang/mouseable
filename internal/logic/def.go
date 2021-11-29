package logic

import (
	"math"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/winapi"
)

type logicState struct {
	fixedSpeed         int
	cursorDX, cursorDY float64
	wheelDX, wheelDY   int
	steppingLogics     []*logicDefinition
	wasCursorMoving    bool
	willDeactivate     bool
}

type logicDefinition struct {
	function *def.FunctionDefinition
	onStart  func(state *logicState)
	onStep   func(state *logicState)
	onStop   func(state *logicState)
}

var logicDefinitions = []*logicDefinition{
	{
		function: def.Deactivate,
		onStop: func(s *logicState) {
			s.willDeactivate = true
		},
	},
	{
		function: def.MoveRight,
		onStep: func(s *logicState) {
			s.cursorDX += cachedDataMap[def.CursorAcceleration].float
		},
	},
	{
		function: def.MoveUp,
		onStep: func(s *logicState) {
			s.cursorDY -= cachedDataMap[def.CursorAcceleration].float
		},
	},
	{
		function: def.MoveLeft,
		onStep: func(s *logicState) {
			s.cursorDX -= cachedDataMap[def.CursorAcceleration].float
		},
	},
	{
		function: def.MoveDown,
		onStep: func(s *logicState) {
			s.cursorDY += cachedDataMap[def.CursorAcceleration].float
		},
	},
	{
		function: def.MoveRightUp,
		onStep: func(s *logicState) {
			spd := cachedDataMap[def.CursorAcceleration].float / 1.41412
			s.cursorDX += spd
			s.cursorDY -= spd
		},
	},
	{
		function: def.MoveLeftUp,
		onStep: func(s *logicState) {
			spd := cachedDataMap[def.CursorAcceleration].float / 1.41412
			s.cursorDX -= spd
			s.cursorDY -= spd
		},
	},
	{
		function: def.MoveRightDown,
		onStep: func(s *logicState) {
			spd := cachedDataMap[def.CursorAcceleration].float / 1.41412
			s.cursorDX += spd
			s.cursorDY += spd
		},
	},
	{
		function: def.MoveLeftDown,
		onStep: func(s *logicState) {
			spd := cachedDataMap[def.CursorAcceleration].float / 1.41412
			s.cursorDX -= spd
			s.cursorDY += spd
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
		onStep: func(s *logicState) {
			s.wheelDY -= cachedDataMap[def.WheelAcceleration].int
		},
	},
	{
		function: def.WheelDown,
		onStep: func(s *logicState) {
			s.wheelDY += cachedDataMap[def.WheelAcceleration].int
		},
	},
	{
		function: def.WheelRight,
		onStep: func(s *logicState) {
			s.wheelDX += cachedDataMap[def.WheelAcceleration].int
		},
	},
	{
		function: def.WheelLeft,
		onStep: func(s *logicState) {
			s.wheelDX -= cachedDataMap[def.WheelAcceleration].int
		},
	},
	{
		function: def.SniperMode,
		onStart: func(s *logicState) {
			s.fixedSpeed = cachedDataMap[def.SniperModeSpeed].int
		},
		onStop: func(s *logicState) {
			s.fixedSpeed = 0
		},
	},
	{
		function: def.TeleportForward,
		onStart: func(s *logicState) {
			if math.Abs(s.cursorDX) < 0.5 && math.Abs(s.cursorDY) < 0.5 {
				return
			}
			distance := cachedDataMap[def.TeleportDistance].int
			var dx int32
			var dy int32
			angle := math.Atan2(s.cursorDX, s.cursorDY)
			dx = int32(math.Round(float64(distance) * math.Sin(angle)))
			dy = int32(math.Round(float64(distance) * math.Cos(angle)))
			winapi.AddCursorPos(dx, dy)
		},
	},
}
