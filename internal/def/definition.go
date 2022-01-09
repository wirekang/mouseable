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
	Activate          = nF("System", "Activate", "Activate Mouseable", Deactivated)
	Deactivate        = nF("System", "Deactivate", "Deactivate Mouseable")
	MoveRight         = nF("Move", "MoveRight", "Move cursor →")
	MoveRightUp       = nF("Move", "MoveRightUp", "Move cursor ↗")
	MoveUp            = nF("Move", "MoveUp", "Move cursor ↑")
	MoveLeftUp        = nF("Move", "MoveLeftUp", "Move cursor ↖")
	MoveLeft          = nF("Move", "MoveLeft", "Move cursor ←")
	MoveLeftDown      = nF("Move", "MoveLeftDown", "Move cursor ↙")
	MoveDown          = nF("Move", "MoveDown", "Move cursor ↓")
	MoveRightDown     = nF("Move", "MoveRightDown", "Move cursor ↘")
	SniperMode        = nF("Move", "SniperMode", "Slow down to increase accuracy")
	SniperModeWheel   = nF("Move", "SniperModeWheel", "Slow down to increase accuracy (Wheel)")
	ClickLeft         = nF("Button", "ClickLeft", "Click left mouse button")
	ClickRight        = nF("Button", "ClickRight", "Click right mouse button")
	ClickMiddle       = nF("Button", "ClickMiddle", "Click middle mouse button")
	WheelUp           = nF("Button", "WheelUp", "Wheel ↑")
	WheelDown         = nF("Button", "WheelDown", "Wheel ↓")
	WheelRight        = nF("Button", "WheelRight", "Wheel →")
	WheelLeft         = nF("Button", "WheelLeft", "Wheel ←")
	TeleportForward   = nF("Teleport", "TeleportForward", "Teleport cursor to the direction it is moving")
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
	DoublePressSpeed      = nD("DoublePressSpeed", "Double press speed in ms", Int)
	CursorAccelerationH   = nD("CursorAccelerationH", "Cursor horizontal acceleration", Float)
	CursorAccelerationV   = nD("CursorAccelerationV", "Cursor vertical acceleration", Float)
	CursorFrictionH       = nD("CursorFrictionH", "Cursor horizontal friction", Float)
	CursorFrictionV       = nD("CursorFrictionV", "Cursor vertical friction", Float)
	WheelAccelerationH    = nD("WheelAccelerationH", "Wheel horizontal acceleration", Int)
	WheelAccelerationV    = nD("WheelAccelerationV", "Wheel vertical acceleration", Int)
	WheelFrictionH        = nD("WheelFrictionH", "Wheel horizontal friction", Int)
	WheelFrictionV        = nD("WheelFrictionV", "Wheel vertical friction", Int)
	SniperModeSpeedH      = nD("SniperModeSpeedH", "Sniper mode horizontal speed", Int)
	SniperModeSpeedV      = nD("SniperModeSpeedV", "Sniper mode vertical speed", Int)
	SniperModeWheelSpeedH = nD("SniperModeWheelSpeedH", "Sniper mode horizontal speed (Wheel)", Int)
	SniperModeWheelSpeedV = nD("SniperModeWheelSpeedV", "Sniper mode vertical speed (Wheel)", Int)
	TeleportDistanceF     = nD("TeleportDistanceF", "TeleportForward distance", Int)
	TeleportDistanceH     = nD("TeleportDistanceH", "Teleport horizontal distance", Int)
	TeleportDistanceV     = nD("TeleportDistanceV", "Teleport vertical distance", Int)
	ShowOverlay           = nD("ShowOverlay", "Show overlay near the cursor when Mouseable activated", Bool)
)
