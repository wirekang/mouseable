package logic

import (
	"math"

	"github.com/wirekang/mouseable/internal/di"
)

func (s *logicState) initCommandTool() {
	s.commandTool = &di.CommandTool{
		Activate: func() {
			s.cmdState.when = di.WhenActivated
			s.overlayManager.Show()
		},
		Deactivate: func() {
			s.cmdState.when = di.WhenDeactivated
			s.overlayManager.Hide()
			s.cursorState.cursorSpeed = emptyVectorInt
			s.cursorState.wheelSpeed = emptyVectorInt
		},
		RegisterCursorAccelerator: func(dir di.Direction) {
			s.cursorState.cursorDirectionMap[dir] = emptyStruct
		},
		UnregisterCursorAccelerator: func(dir di.Direction) {
			delete(s.cursorState.cursorDirectionMap, dir)
			if len(s.cursorState.cursorDirectionMap) == 0 {
				s.cursorState.cursorSpeed = emptyVectorInt
			}
		},
		RegisterWheelAccelerator: func(dir di.Direction) {
			s.cursorState.wheelDirectionMap[dir] = emptyStruct

		},
		UnregisterWheelAccelerator: func(dir di.Direction) {
			delete(s.cursorState.wheelDirectionMap, dir)
			if len(s.cursorState.wheelDirectionMap) == 0 {
				s.cursorState.wheelSpeed = emptyVectorInt
			}
		},
		FixCursorSpeed: func() {
			s.cursorState.maxCursorSpeed = s.configCache.cursorSniperSpeed
		},
		UnfixCursorSpeed: func() {
			s.cursorState.maxCursorSpeed = s.configCache.cursorMaxSpeed
		},
		FixWheelSpeed: func() {
			s.cursorState.maxWheelSpeed = s.configCache.wheelSniperSpeed
		},
		UnfixWheelSpeed: func() {
			s.cursorState.maxWheelSpeed = s.configCache.wheelMaxSpeed
		},
		MouseDown: func(button di.MouseButton) {
			go s.hookManager.MouseDown(button)
		},
		MouseUp: func(button di.MouseButton) {
			go s.hookManager.MouseUp(button)
		},
		Teleport: func(dir di.Direction) {
			s.channel.cursorBuffer <- directionToVectorInt(dir, s.configCache.teleportDistance)
		},
		TeleportForward: func() {
			if math.Abs(float64(s.cursorState.cursorSpeed.x)) > 0.3 ||
				math.Abs(float64(s.cursorState.cursorSpeed.y)) > 0.3 {
				distance := s.configCache.teleportDistance
				angle := math.Atan2(
					float64(s.cursorState.cursorSpeed.x),
					float64(s.cursorState.cursorSpeed.y),
				)
				s.cursorState.lastTeleportForward = vectorInt{
					x: int(math.Round(float64(distance) * math.Sin(angle))),
					y: int(math.Round(float64(distance) * math.Cos(angle))),
				}
			}
			s.channel.cursorBuffer <- s.cursorState.lastTeleportForward
		},
		Toggle: func() {
			if s.cmdState.when == di.WhenDeactivated {
				s.commandTool.Activate()
			} else {
				s.commandTool.Deactivate()
			}
		},
	}
}

const slow = 1 / math.Sqrt2

var directionVectorMap = map[di.Direction]vectorFloat{
	di.DirectionRight:     {x: 1, y: 0},
	di.DirectionRightUp:   {x: slow, y: -slow},
	di.DirectionUp:        {x: 0, y: -1},
	di.DirectionLeftUp:    {x: -slow, y: -slow},
	di.DirectionLeft:      {x: -1, y: 0},
	di.DirectionLeftDown:  {x: -slow, y: slow},
	di.DirectionDown:      {x: 0, y: 1},
	di.DirectionRightDown: {x: slow, y: slow},
}

func directionToVectorInt(d di.Direction, distance int) (r vectorInt) {
	f := directionVectorMap[d]
	r.x = int(math.Round(f.x * float64(distance)))
	r.y = int(math.Round(f.y * float64(distance)))
	return
}
