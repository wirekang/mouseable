// Package def declares types for prevent circular dependencies.
package def

import (
	"github.com/wirekang/mouseable/internal/typ"
)

func New() typ.DefinitionManager {
	m := &manager{
		cmdOrderMap:        map[typ.CommandName]int{},
		cmdWhenMap:         map[typ.CommandName]typ.When{},
		cmdDescriptionMap:  map[typ.CommandName]string{},
		dataTypeMap:        map[typ.DataName]typ.DataType{},
		dataDescriptionMap: map[typ.DataName]string{},
		nextFuncOrder:      0,
	}
	m.nc("Activate", "Activate Mouseable", typ.Deactivated)
	m.nc("Deactivate", "Deactivate Mouseable", typ.Activated)
	m.nc("MoveRight", "Move cursor →", typ.Activated)
	m.nc("MoveRightUp", "Move cursor ↗", typ.Activated)
	m.nc("MoveUp", "Move cursor ↑", typ.Activated)
	m.nc("MoveLeftUp", "Move cursor ↖", typ.Activated)
	m.nc("MoveLeft", "Move cursor ←", typ.Activated)
	m.nc("MoveLeftDown", "Move cursor ↙", typ.Activated)
	m.nc("MoveDown", "Move cursor ↓", typ.Activated)
	m.nc("MoveRightDown", "Move cursor ↘", typ.Activated)
	m.nc("SniperMode", "Slow down to increase accuracy", typ.Activated)
	m.nc("SniperModeWheel", "Slow down to increase accuracy (Wheel)", typ.Activated)
	m.nc("ClickLeft", "Click left mouse button", typ.Activated)
	m.nc("ClickRight", "Click right mouse button", typ.Activated)
	m.nc("ClickMiddle", "Click middle mouse button", typ.Activated)
	m.nc("WheelUp", "Wheel ↑", typ.Activated)
	m.nc("WheelDown", "Wheel ↓", typ.Activated)
	m.nc("WheelRight", "Wheel →", typ.Activated)
	m.nc("WheelLeft", "Wheel ←", typ.Activated)
	m.nc("TeleportForward", "Teleport cursor to the direction it is moving", typ.Activated)
	m.nc("TeleportRight", "Teleport cursor →", typ.Activated)
	m.nc("TeleportRightUp", "Teleport cursor ↗", typ.Activated)
	m.nc("TeleportUp", "Teleport cursor ↑", typ.Activated)
	m.nc("TeleportLeftUp", "Teleport cursor ↖", typ.Activated)
	m.nc("TeleportLeft", "Teleport cursor ←", typ.Activated)
	m.nc("TeleportLeftDown", "Teleport cursor ↙", typ.Activated)
	m.nc("TeleportDown", "Teleport cursor ↓", typ.Activated)
	m.nc("TeleportRightDown", "Teleport cursor ↘", typ.Activated)

	m.nd("DoublePressSpeed", "Double press speed in ms", typ.Int)
	m.nd("CursorAccelerationH", "Cursor horizontal acceleration", typ.Float)
	m.nd("CursorAccelerationV", "Cursor vertical acceleration", typ.Float)
	m.nd("CursorFrictionH", "Cursor horizontal friction", typ.Float)
	m.nd("CursorFrictionV", "Cursor vertical friction", typ.Float)
	m.nd("WheelAccelerationH", "Wheel horizontal acceleration", typ.Int)
	m.nd("WheelAccelerationV", "Wheel vertical acceleration", typ.Int)
	m.nd("WheelFrictionH", "Wheel horizontal friction", typ.Int)
	m.nd("WheelFrictionV", "Wheel vertical friction", typ.Int)
	m.nd("SniperModeSpeedH", "Sniper mode horizontal speed", typ.Int)
	m.nd("SniperModeSpeedV", "Sniper mode vertical speed", typ.Int)
	m.nd("SniperModeWheelSpeedH", "Sniper mode horizontal speed (Wheel)", typ.Int)
	m.nd("SniperModeWheelSpeedV", "Sniper mode vertical speed (Wheel)", typ.Int)
	m.nd("TeleportDistanceF", "TeleportForward distance", typ.Int)
	m.nd("TeleportDistanceH", "Teleport horizontal distance", typ.Int)
	m.nd("TeleportDistanceV", "Teleport vertical distance", typ.Int)
	m.nd("ShowOverlay", "Show overlay when Mouseable activated", typ.Bool)

	return m
}
