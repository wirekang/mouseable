package logic

import (
	"math"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/winapi"
)

type logicState struct {
	fixedSpeedH        int
	fixedSpeedV        int
	cursorDX, cursorDY float64
	wheelDX, wheelDY   int
	willActivate       bool
	willDeactivate     bool
}

type keyState struct {
	modInfo  modInfo
	keyCodes []uint32
}

type modInfo struct {
	isWin,
	isControl,
	isAlt,
	isShift bool
}

type logicDefinition struct {
	function *def.FunctionDefinition
	onStart  func(state *logicState)
	onStep   func(state *logicState)
	onStop   func(state *logicState)
}

var logicDefinitions = []*logicDefinition{
	{
		function: def.Activate,
		onStart: func(s *logicState) {
			s.willActivate = true
		},
	},
	{
		function: def.Deactivate,
		onStop: func(s *logicState) {
			s.willDeactivate = true
		},
	},
	{
		function: def.MoveRight,
		onStep: func(s *logicState) {
			s.cursorDX += cachedDataMap[def.CursorAccelerationH].float
		},
	},
	{
		function: def.MoveUp,
		onStep: func(s *logicState) {
			s.cursorDY -= cachedDataMap[def.CursorAccelerationV].float
		},
	},
	{
		function: def.MoveLeft,
		onStep: func(s *logicState) {
			s.cursorDX -= cachedDataMap[def.CursorAccelerationH].float
		},
	},
	{
		function: def.MoveDown,
		onStep: func(s *logicState) {
			s.cursorDY += cachedDataMap[def.CursorAccelerationV].float
		},
	},
	{
		function: def.MoveRightUp,
		onStep: func(s *logicState) {
			hs := cachedDataMap[def.CursorAccelerationH].float / 1.4
			vs := cachedDataMap[def.CursorAccelerationV].float / 1.4
			s.cursorDX += hs
			s.cursorDY -= vs
		},
	},
	{
		function: def.MoveLeftUp,
		onStep: func(s *logicState) {
			hs := cachedDataMap[def.CursorAccelerationH].float / 1.4
			vs := cachedDataMap[def.CursorAccelerationV].float / 1.4
			s.cursorDX -= hs
			s.cursorDY -= vs
		},
	},
	{
		function: def.MoveRightDown,
		onStep: func(s *logicState) {
			hs := cachedDataMap[def.CursorAccelerationH].float / 1.4
			vs := cachedDataMap[def.CursorAccelerationV].float / 1.4
			s.cursorDX += hs
			s.cursorDY += vs
		},
	},
	{
		function: def.MoveLeftDown,
		onStep: func(s *logicState) {
			hs := cachedDataMap[def.CursorAccelerationH].float / 1.4
			vs := cachedDataMap[def.CursorAccelerationV].float / 1.4
			s.cursorDX -= hs
			s.cursorDY += vs
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
			s.wheelDY += cachedDataMap[def.WheelAccelerationV].int
		},
	},
	{
		function: def.WheelDown,
		onStep: func(s *logicState) {
			s.wheelDY -= cachedDataMap[def.WheelAccelerationV].int
		},
	},
	{
		function: def.WheelRight,
		onStep: func(s *logicState) {
			s.wheelDX -= cachedDataMap[def.WheelAccelerationH].int
		},
	},
	{
		function: def.WheelLeft,
		onStep: func(s *logicState) {
			s.wheelDX += cachedDataMap[def.WheelAccelerationH].int
		},
	},
	{
		function: def.SniperMode,
		onStart: func(s *logicState) {
			s.fixedSpeedH = cachedDataMap[def.SniperModeSpeedH].int
			s.fixedSpeedV = cachedDataMap[def.SniperModeSpeedV].int
		},
		onStop: func(s *logicState) {
			s.fixedSpeedH = 0
			s.fixedSpeedV = 0
		},
	},
	{
		function: def.TeleportForward,
		onStart: func(s *logicState) {
			if math.Abs(s.cursorDX) < 0.5 && math.Abs(s.cursorDY) < 0.5 {
				return
			}
			distance := cachedDataMap[def.TeleportDistanceF].int
			var dx int
			var dy int
			angle := math.Atan2(s.cursorDX, s.cursorDY)
			dx = int(math.Round(float64(distance) * math.Sin(angle)))
			dy = int(math.Round(float64(distance) * math.Cos(angle)))
			winapi.AddCursorPos(dx, dy)
		},
	},
}
