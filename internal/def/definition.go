package def

import (
	"github.com/wirekang/mouseable/internal/di"
)

func New() di.DefinitionManager {
	m := &manager{
		cmdDefMap:  map[di.CommandName]*commandDef{},
		dataDefMap: map[di.DataName]*dataDef{},
	}
	m.insertCommand(
		"activate",
		"Activate Mouseable",
		di.WhenDeactivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Activate()
			},
		},
	)
	m.insertCommand(
		"activate-temp",
		"Activate Mouseable temporarily",
		di.WhenDeactivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Activate()
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.Deactivate()
			},
		},
	)
	m.insertCommand(
		"deactivate",
		"Deactivate Mouseable",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Deactivate()
			},
		},
	)
	m.insertCommand(
		"toggle", "Toggle Activate <-> Deactivate", di.WhenAnytime, &di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Toggle()
			},
		},
	)
	m.insertCommand(
		"deactivate-temp",
		"Deactivate Mouseable temporarily",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Deactivate()
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.Activate()
			},
		},
	)
	m.insertCommand(
		"move-right",
		"Move cursor →",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterCursorAccelerator(di.DirectionRight)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterCursorAccelerator(di.DirectionRight)
			},
		},
	)
	m.insertCommand(
		"move-right-up",
		"Move cursor ↗",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterCursorAccelerator(di.DirectionRight | di.DirectionUp)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterCursorAccelerator(di.DirectionRight | di.DirectionUp)
			},
		},
	)
	m.insertCommand(
		"move-up",
		"Move cursor ↑",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterCursorAccelerator(di.DirectionUp)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterCursorAccelerator(di.DirectionUp)
			},
		},
	)
	m.insertCommand(
		"move-left-up",
		"Move cursor ↖",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterCursorAccelerator(di.DirectionLeft | di.DirectionUp)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterCursorAccelerator(di.DirectionLeft | di.DirectionUp)
			},
		},
	)
	m.insertCommand(
		"move-left",
		"Move cursor ←",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterCursorAccelerator(di.DirectionLeft)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterCursorAccelerator(di.DirectionLeft)
			},
		},
	)
	m.insertCommand(
		"move-left-down",
		"Move cursor ↙",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterCursorAccelerator(di.DirectionLeft | di.DirectionDown)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterCursorAccelerator(di.DirectionLeft | di.DirectionDown)
			},
		},
	)
	m.insertCommand(
		"move-down",
		"Move cursor ↓",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterCursorAccelerator(di.DirectionDown)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterCursorAccelerator(di.DirectionDown)
			},
		},
	)
	m.insertCommand(
		"move-right-down",
		"Move cursor ↘",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterCursorAccelerator(di.DirectionRight | di.DirectionDown)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterCursorAccelerator(di.DirectionRight | di.DirectionDown)
			},
		},
	)
	m.insertCommand(
		"sniper-mode",
		"Slow down to increase accuracy",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.FixCursorSpeed()
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnfixCursorSpeed()
			},
		},
	)
	m.insertCommand(
		"sniper-mode-wheel",
		"Slow down to increase accuracy (MouseWheel)",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.FixWheelSpeed()
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnfixWheelSpeed()
			},
		},
	)
	m.insertCommand(
		"click-left",
		"Click left mouse button",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.MouseDown(di.ButtonLeft)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.MouseUp(di.ButtonLeft)
			},
		},
	)
	m.insertCommand(
		"click-right",
		"Click right mouse button",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.MouseDown(di.ButtonRight)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.MouseUp(di.ButtonRight)
			},
		},
	)
	m.insertCommand(
		"click-middle",
		"Click middle mouse button",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.MouseDown(di.ButtonMiddle)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.MouseUp(di.ButtonMiddle)
			},
		},
	)
	m.insertCommand(
		"wheel-up",
		"MouseWheel ↑",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterWheelAccelerator(di.DirectionUp)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterWheelAccelerator(di.DirectionUp)
			},
		},
	)
	m.insertCommand(
		"wheel-down",
		"MouseWheel ↓",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterWheelAccelerator(di.DirectionDown)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterWheelAccelerator(di.DirectionDown)
			},
		},
	)
	m.insertCommand(
		"wheel-right",
		"MouseWheel →",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterWheelAccelerator(di.DirectionRight)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterWheelAccelerator(di.DirectionRight)
			},
		},
	)
	m.insertCommand(
		"wheel-left",
		"MouseWheel ←",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.RegisterWheelAccelerator(di.DirectionLeft)
			},
			OnEnd: func(tool *di.CommandTool) {
				tool.UnregisterWheelAccelerator(di.DirectionLeft)
			},
		},
	)
	m.insertCommand(
		"teleport-forward",
		"Teleport cursor to the direction it is moving",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.TeleportForward()
			},
		},
	)
	m.insertCommand(
		"teleport-right",
		"Teleport cursor →",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Teleport(di.DirectionRight)
			},
		},
	)
	m.insertCommand(
		"teleport-right-up",
		"Teleport cursor ↗",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Teleport(di.DirectionRight | di.DirectionUp)
			},
		},
	)
	m.insertCommand(
		"teleport-up",
		"Teleport cursor ↑",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Teleport(di.DirectionUp)
			},
		},
	)
	m.insertCommand(
		"teleport-left-up",
		"Teleport cursor ↖",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Teleport(di.DirectionLeft | di.DirectionUp)
			},
		},
	)
	m.insertCommand(
		"teleport-left",
		"Teleport cursor ←",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Teleport(di.DirectionLeft)
			},
		},
	)
	m.insertCommand(
		"teleport-left-down",
		"Teleport cursor ↙",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Teleport(di.DirectionLeft | di.DirectionDown)
			},
		},
	)
	m.insertCommand(
		"teleport-down",
		"Teleport cursor ↓",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Teleport(di.DirectionDown)
			},
		},
	)
	m.insertCommand(
		"teleport-right-down",
		"Teleport cursor ↘",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Teleport(di.DirectionRight | di.DirectionDown)
			},
		},
	)
	m.insertCommand(
		"attach-right",
		"Attach cursor rightmost.",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Attach(di.DirectionRight)
			},
		},
	)
	m.insertCommand(
		"attach-left",
		"Attach cursor leftmost.",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Attach(di.DirectionLeft)
			},
		},
	)
	m.insertCommand(
		"attach-up",
		"Attach cursor uppermost.",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Attach(di.DirectionUp)
			},
		},
	)
	m.insertCommand(
		"attach-down",
		"Attach cursor bottommost.",
		di.WhenActivated,
		&di.Command{
			OnBegin: func(tool *di.CommandTool) {
				tool.Attach(di.DirectionDown)
			},
		},
	)

	m.insertData("key-timeout", "Key press timeout for continuous input in ms", di.Int, 150)
	m.insertData("cursor-acceleration", "Cursor acceleration", di.Float, 0.5)
	m.insertData("cursor-max-speed", "Cursor max speed", di.Int, 10)
	m.insertData("wheel-acceleration", "Wheel acceleration", di.Float, 4.0)
	m.insertData("wheel-max-speed", "Wheel max speed", di.Int, 40)
	m.insertData("cursor-sniper-speed", "Sniper mode speed", di.Int, 2)
	m.insertData("wheel-sniper-speed", "Sniper mode speed (Wheel)", di.Int, 8)
	m.insertData("teleport-distance", "Teleport distance", di.Int, 300)
	m.insertData("show-overlay", "Show overlay when Mouseable activated", di.Bool, true)
	m.insertData(
		"cursor-factor",
		"If value is 1, cursor moves at the same vertical/horizontal speed.\n\n"+
			" If value is 2, horizontal speed of cursor is vertical speed x 2.\n\n"+
			"If value is 0.5,horizontal speed of cursor is vertical speed x 0.5",
		di.Float, 1.0,
	)
	m.insertData("wheel-factor", "Same as cursor-factor", di.Float, 1.0)
	m.insertData("teleport-factor", "Same as cursor-factor", di.Float, 1.0)
	m.insertData("fast-diagonals", "If true, cursor moves faster on the diagonal.", di.Bool, false)

	return m
}
