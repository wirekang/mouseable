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
	m.nc("activate", "Activate Mouseable", typ.Deactivated)
	m.nc("deactivate", "Deactivate Mouseable", typ.Activated)
	m.nc("move-right", "Move cursor →", typ.Activated)
	m.nc("move-right-up", "Move cursor ↗", typ.Activated)
	m.nc("move-up", "Move cursor ↑", typ.Activated)
	m.nc("move-left-up", "Move cursor ↖", typ.Activated)
	m.nc("move-left", "Move cursor ←", typ.Activated)
	m.nc("move-left-down", "Move cursor ↙", typ.Activated)
	m.nc("move-down", "Move cursor ↓", typ.Activated)
	m.nc("move-right-down", "Move cursor ↘", typ.Activated)
	m.nc("sniper-mode", "Slow down to increase accuracy", typ.Activated)
	m.nc("sniper-mode-wheel", "Slow down to increase accuracy (Wheel)", typ.Activated)
	m.nc("click-left", "Click left mouse button", typ.Activated)
	m.nc("click-right", "Click right mouse button", typ.Activated)
	m.nc("click-middle", "Click middle mouse button", typ.Activated)
	m.nc("wheel-up", "Wheel ↑", typ.Activated)
	m.nc("wheel-down", "Wheel ↓", typ.Activated)
	m.nc("wheel-right", "Wheel →", typ.Activated)
	m.nc("wheel-left", "Wheel ←", typ.Activated)
	m.nc("teleport-forward", "Teleport cursor to the direction it is moving", typ.Activated)
	m.nc("teleport-right", "Teleport cursor →", typ.Activated)
	m.nc("teleport-right-up", "Teleport cursor ↗", typ.Activated)
	m.nc("teleport-up", "Teleport cursor ↑", typ.Activated)
	m.nc("teleport-left-up", "Teleport cursor ↖", typ.Activated)
	m.nc("teleport-left", "Teleport cursor ←", typ.Activated)
	m.nc("teleport-left-down", "Teleport cursor ↙", typ.Activated)
	m.nc("teleport-down", "Teleport cursor ↓", typ.Activated)
	m.nc("teleport-right-down", "Teleport cursor ↘", typ.Activated)

	m.nd("double-press-speed", "Double press speed in ms", typ.Int)
	m.nd("cursor-acceleration-h", "Cursor horizontal acceleration", typ.Float)
	m.nd("cursor-acceleration-v", "Cursor vertical acceleration", typ.Float)
	m.nd("cursor-friction-h", "Cursor horizontal friction", typ.Float)
	m.nd("cursor-friction-v", "Cursor vertical friction", typ.Float)
	m.nd("wheel-acceleration-h", "Wheel horizontal acceleration", typ.Int)
	m.nd("wheel-acceleration-v", "Wheel vertical acceleration", typ.Int)
	m.nd("wheel-friction-h", "Wheel horizontal friction", typ.Int)
	m.nd("wheel-friction-v", "Wheel vertical friction", typ.Int)
	m.nd("sniper-mode-speed-h", "Sniper mode horizontal speed", typ.Int)
	m.nd("sniper-mode-speed-v", "Sniper mode vertical speed", typ.Int)
	m.nd("sniper-mode-wheel-speed-h", "Sniper mode horizontal speed (Wheel)", typ.Int)
	m.nd("sniper-mode-wheel-speed-v", "Sniper mode vertical speed (Wheel)", typ.Int)
	m.nd("teleport-distance-f", "TeleportForward distance", typ.Int)
	m.nd("teleport-distance-h", "Teleport horizontal distance", typ.Int)
	m.nd("teleport-distance-v", "Teleport vertical distance", typ.Int)
	m.nd("show-overlay", "Show overlay when Mouseable activated", typ.Bool)

	return m
}
