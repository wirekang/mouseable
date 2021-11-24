package di

import (
	"github.com/wirekang/mouseable/internal/data"
	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/hook"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/logic"
	"github.com/wirekang/mouseable/internal/overlay"
	"github.com/wirekang/mouseable/internal/view"
)

func Init() {
	logic.DI.SetCursorPos = func(x, y int) {
		hook.SetCursorPos(x, y)
		lg.Logf("logic.DI.SetCursorPos(%d, %d)", x, y)
	}
	logic.DI.GetCursorPos = func() (x, y int) {
		x, y = hook.GetCursorPos()
		lg.Logf("logic.DI.GetCursorPos() (%d, %d)", x, y)
		return
	}
	logic.DI.AddCursorPos = func(dx, dy int32) {
		hook.AddCursorPos(dx, dy)
		lg.Logf("logic.DI.AddCursorPos(%d, %d)", dx, dy)
	}
	logic.DI.MouseDown = func(button int) {
		hook.MouseDown(button)
		lg.Logf("logic.DI.MouseDown(%d)", button)
	}
	logic.DI.MouseUp = func(button int) {
		hook.MouseUp(button)
		lg.Logf("logic.DI.MouseUp(%d)", button)
	}
	logic.DI.Wheel = func(amount int, isHorizontal bool) {
		hook.Wheel(amount, isHorizontal)
		lg.Logf("logic.DI.Whell(%d, %v)", amount, isHorizontal)
	}
	logic.DI.OnCursorMove = func() {
		overlay.OnCursorMove()
		lg.Logf("logic.DI.OnCursorMove()")
	}
	logic.DI.OnCursorStop = func() {
		overlay.OnCursorStop()
		lg.Logf("logic.DI.OnCursorStop()")
	}
	logic.DI.Unhook = func() {
		hook.Unhook()
		lg.Logf("logic.DI.Unhook()")
	}
	hook.DI.OnKey = func(keyCode uint32, isDown bool) (preventDefault bool) {
		preventDefault = logic.OnKey(keyCode, isDown)
		lg.Logf("hook.DI.OnKey(%d, %v) %v", keyCode, isDown, preventDefault)
		return
	}
	hook.DI.OnHook = func() {
		overlay.OnHook()
		lg.Logf("hook.DI.OnHook()")
	}
	hook.DI.OnUnhook = func() {
		overlay.OnUnhook()
		lg.Logf("hook.DI.OnUnhook()")
	}
	view.DI.LoadConfig = func() (config def.Config) {
		config = data.LoadConfig()
		lg.Logf("view.DI.LoadConfig() %+v", config)
		return
	}
	view.DI.SaveConfig = func(config def.Config) {
		data.SaveConfig(config)
		lg.Logf("view.DI.SaveConfig(%+v)", config)
	}
	data.DI.SetConfig = func(config def.Config) {
		hook.SetKey(config)
		logic.SetConfig(config)
		lg.Logf("data.DI.SetConfig(%+v)", config)
	}
}
