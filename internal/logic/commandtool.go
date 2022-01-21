package logic

import (
	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/logic/mover"
)

func (s *logicState) initCommandTool() {
	s.commandTool = &di.CommandTool{
		Activate: func() {
			s.cmdState.when = di.WhenActivated
			s.overlayManager.Show()
			s.uiManager.SetTrayIconEnabled(true)
		},
		Deactivate: func() {
			s.cmdState.when = di.WhenDeactivated
			s.overlayManager.Hide()
			s.uiManager.SetTrayIconEnabled(false)
			s.cursorState.cursorMover.SetSpeed(0)
		},
		RegisterCursorAccelerator: func(dir di.Direction) {
			s.cursorState.cursorMover.AddDirection(dir)
		},
		UnregisterCursorAccelerator: func(dir di.Direction) {
			s.cursorState.cursorMover.RemoveDirection(dir)
		},
		RegisterWheelAccelerator: func(dir di.Direction) {
			s.cursorState.wheelMover.AddDirection(dir)
		},
		UnregisterWheelAccelerator: func(dir di.Direction) {
			s.cursorState.wheelMover.RemoveDirection(dir)
		},
		FixCursorSpeed: func() {
			s.cursorState.cursorMover.SetMaxSpeed(s.configCache.cursorSniperSpeed)
		},
		UnfixCursorSpeed: func() {
			s.cursorState.cursorMover.SetMaxSpeed(s.configCache.cursorMaxSpeed)
		},
		FixWheelSpeed: func() {
			s.cursorState.wheelMover.SetMaxSpeed(s.configCache.wheelSniperSpeed)
		},
		UnfixWheelSpeed: func() {
			s.cursorState.wheelMover.SetMaxSpeed(s.configCache.wheelMaxSpeed)
		},
		MouseDown: func(button di.MouseButton) {
			go s.hookManager.MouseDown(button)
		},
		MouseUp: func(button di.MouseButton) {
			go s.hookManager.MouseUp(button)
		},
		Teleport: func(dir di.Direction) {
			s.cursorState.teleportMover.SetDirection(dir)
			s.channel.cursorBuffer <- s.cursorState.teleportMover.Vector()
		},
		TeleportForward: func() {
			s.cursorState.teleportMover.SetDirection(s.cursorState.cursorMover.Direction())
			s.channel.cursorBuffer <- s.cursorState.teleportMover.Vector()
		},
		Toggle: func() {
			if s.cmdState.when == di.WhenDeactivated {
				s.commandTool.Activate()
			} else {
				s.commandTool.Deactivate()
			}
		},
		Attach: func(dir di.Direction) {
			m := mover.Mover{}
			m.SetDirection(dir)
			m.SetMaxSpeed(20000)
			m.SetSpeed(20000)
			s.channel.cursorBuffer <- m.Vector()
		},
	}
}
