// Package def declares types for prevent circular dependencies.
package def

import (
	"fmt"
)

type FunctionMap map[*FunctionDefinition]FunctionKey
type DataMap map[*DataDefinition]DataValue

type Config struct {
	FunctionMap FunctionMap
	DataMap     DataMap
}

var FunctionDefinitions []*FunctionDefinition
var DataDefinitions []*DataDefinition

var nextOrder = 0

func nF(category, name, desc string, when ...When) (f *FunctionDefinition) {
	nextOrder++
	f = new(FunctionDefinition)
	f.Order = nextOrder
	f.Category = category
	f.Name = name
	f.Description = desc
	if len(when) == 0 {
		f.When = Activated
	} else {
		f.When = when[0]
	}
	FunctionDefinitions = append(FunctionDefinitions, f)
	return
}

func nD(name, desc string, t Type) (d *DataDefinition) {
	d = new(DataDefinition)
	d.Name = name
	d.Description = desc
	d.Type = t
	DataDefinitions = append(DataDefinitions, d)
	return
}

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
		"Move-2",
		"TeleportForward", "Teleport cursor to the direction it is moving",
	)
	TeleportRight     = nF("Move-2", "TeleportRight", "Teleport cursor →")
	TeleportRightUp   = nF("Move-2", "TeleportRightUp", "Teleport cursor ↗")
	TeleportUp        = nF("Move-2", "TeleportUp", "Teleport cursor ↑")
	TeleportLeftUp    = nF("Move-2", "TeleportLeftUp", "Teleport cursor ↖")
	TeleportLeft      = nF("Move-2", "TeleportLeft", "Teleport cursor ←")
	TeleportLeftDown  = nF("Move-2", "TeleportLeftDown", "Teleport cursor ↙")
	TeleportDown      = nF("Move-2", "TeleportDown", "Teleport cursor ↓")
	TeleportRightDown = nF("Move-2", "TeleportRightDown", "Teleport cursor ↘")
)

var (
	CursorAcceleration = nD("CursorAcceleration", "Cursor acceleration Value", Float)
	CursorFriction     = nD("CursorFriction", "Cursor friction Value", Float)
	WheelAcceleration  = nD("WheelAcceleration", "Wheel acceleration Value", Int)
	WheelFriction      = nD("WheelFriction", "Wheel friction Value", Int)
	SniperModeSpeed    = nD("SniperModeSpeed", "Speed in sniper mode", Int)
	TeleportDistance   = nD("TeleportDistance", "Teleport distance", Int)
)

func GetFunctionDefinitionByName(name string) (f *FunctionDefinition) {
	for i := range FunctionDefinitions {
		if FunctionDefinitions[i].Name == name {
			return FunctionDefinitions[i]
		}
	}
	panic(fmt.Sprintf("%s is not a FunctionDifinition name", name))
}

func GetDataDefinitionByName(name string) (d *DataDefinition) {
	for i := range DataDefinitions {
		if DataDefinitions[i].Name == name {
			return DataDefinitions[i]
		}
	}
	panic(fmt.Sprintf("%s is not a DataDifinition name", name))
}
