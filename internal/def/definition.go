// Package def declares types for prevent circular dependencies.
package def

type FunctionMap map[*FunctionDefinition]FunctionKey
type DataMap map[*DataDefinition]DataValue

type Config struct {
	FunctionMap FunctionMap
	DataMap     DataMap
}

var FunctionDefinitions []*FunctionDefinition
var DataDefinitions []*DataDefinition
var FunctionNameMap = map[string]*FunctionDefinition{}
var DataNameMap = map[string]*DataDefinition{}

var (
	Activate = nF(
		"System", "Activate", "Activate Mouseable", Deactivated,
	)
	Deactivate      = nF("System", "Deactivate", "Deactivate Mouseable")
	MoveRight       = nF("Move", "MoveRight", "Move cursor →")
	MoveRightUp     = nF("Move", "MoveRightUp", "Move cursor ↗")
	MoveUp          = nF("Move", "MoveUp", "Move cursor ↑")
	MoveLeftUp      = nF("Move", "MoveLeftUp", "Move cursor ↖")
	MoveLeft        = nF("Move", "MoveLeft", "Move cursor ←")
	MoveLeftDown    = nF("Move", "MoveLeftDown", "Move cursor ↙")
	MoveDown        = nF("Move", "MoveDown", "Move cursor ↓")
	MoveRightDown   = nF("Move", "MoveRightDown", "Move cursor ↘")
	SniperMode      = nF("Move", "SniperMode", "Slow down to increase accuracy")
	ClickLeft       = nF("Button", "ClickLeft", "Click left mouse button")
	ClickRight      = nF("Button", "ClickRight", "Click right mouse button")
	ClickMiddle     = nF("Button", "ClickMiddle", "Click middle mouse button")
	WheelUp         = nF("Button", "WheelUp", "Wheel ↑")
	WheelDown       = nF("Button", "WheelDown", "Wheel ↓")
	WheelRight      = nF("Button", "WheelRight", "Wheel →")
	WheelLeft       = nF("Button", "WheelLeft", "Wheel ←")
	TeleportForward = nF(
		"Teleport",
		"TeleportForward", "Teleport cursor to the direction it is moving",
	)
	TeleportRight     = nF("Teleport", "TeleportRight", "Teleport cursor →")
	TeleportRightUp   = nF("Teleport", "TeleportRightUp", "Teleport cursor ↗")
	TeleportUp        = nF("Teleport", "TeleportUp", "Teleport cursor ↑")
	TeleportLeftUp    = nF("Teleport", "TeleportLeftUp", "Teleport cursor ↖")
	TeleportLeft      = nF("Teleport", "TeleportLeft", "Teleport cursor ←")
	TeleportLeftDown  = nF("Teleport", "TeleportLeftDown", "Teleport cursor ↙")
	TeleportDown      = nF("Teleport", "TeleportDown", "Teleport cursor ↓")
	TeleportRightDown = nF("Teleport", "TeleportRightDown", "Teleport cursor ↘")
)

var (
	CursorAcceleration = nD("CursorAcceleration", "Cursor acceleration Value", Float)
	CursorFriction     = nD("CursorFriction", "Cursor friction Value", Float)
	WheelAcceleration  = nD("WheelAcceleration", "Wheel acceleration Value", Int)
	WheelFriction      = nD("WheelFriction", "Wheel friction Value", Int)
	SniperModeSpeed    = nD("SniperModeSpeed", "Speed in sniper mode", Int)
	TeleportDistance   = nD("TeleportDistance", "Teleport distance", Int)
)
