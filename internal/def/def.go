// Package def declares types for prevent circular dependencies.
package def

type def struct {
	Order       int
	Name        string
	Description string
}

type FunctionDef def

type DataDef def

type HotKeyDef def

type HotKey struct {
	IsAlt     bool
	IsControl bool
	IsShift   bool
	IsWin     bool
	KeyCode   uint32
}

type Config struct {
	HotKeyMap          map[*HotKeyDef]HotKey
	FunctionKeyCodeMap map[*FunctionDef]uint32
	DataValueMap       map[*DataDef]float64
}

var HotKeyNameMap = map[string]*HotKeyDef{}
var FunctionNameMap = map[string]*FunctionDef{}
var DataNameMap = map[string]*DataDef{}
var HotKeyDefs []*HotKeyDef
var FunctionDefs []*FunctionDef
var DataDefs []*DataDef

var nextHotKeyOrder = 0
var nextFunctionOrder = 0
var nextDataOrder = 0

func newH(name, des string) *HotKeyDef {
	h := &HotKeyDef{Order: nextHotKeyOrder, Name: name, Description: des}
	nextHotKeyOrder++
	HotKeyNameMap[name] = h
	HotKeyDefs = append(HotKeyDefs, h)
	return h
}

func newF(name, des string) *FunctionDef {
	f := &FunctionDef{Order: nextFunctionOrder, Name: name, Description: des}
	nextFunctionOrder++
	FunctionNameMap[name] = f
	FunctionDefs = append(FunctionDefs, f)
	return f
}

func newD(name, des string) *DataDef {
	d := &DataDef{Order: nextDataOrder, Name: name, Description: des}
	nextDataOrder++
	DataNameMap[name] = d
	DataDefs = append(DataDefs, d)
	return d
}

var (
	Activate = newH("Activate", "Activate Mouseable")
)

var (
	Deactivate    = newF("Deactivate", "Deactivate Mouseable")
	MoveRight     = newF("MoveRight", "Move cursor →")
	MoveRightUp   = newF("MoveRightUp", "Move cursor ↗")
	MoveUp        = newF("MoveUp", "Move cursor ↑")
	MoveLeftUp    = newF("MoveLeftUp", "Move cursor ↖")
	MoveLeft      = newF("MoveLeft", "Move cursor ←")
	MoveLeftDown  = newF("MoveLeftDown", "Move cursor ↙")
	MoveDown      = newF("MoveDown", "Move cursor ↓")
	MoveRightDown = newF("MoveRightDown", "Move cursor ↘")
	ClickLeft     = newF("ClickLeft", "Click left mouse button")
	ClickRight    = newF("ClickRight", "Click right mouse button")
	ClickMiddle   = newF("ClickMiddle", "Click middle mouse button")
	WheelUp       = newF("WheelUp", "Wheel ↑")
	WheelDown     = newF("WheelDown", "Wheel ↓")
	WheelRight    = newF("WheelRight", "Wheel →")
	WheelLeft     = newF("WheelLeft", "Wheel ←")
	SniperMode    = newF("SniperMode", "Slow down to increase accuracy")
	Flash         = newF(
		"Flash", "Teleport cursor to the direction it is moving",
	)
)

var (
	Acceleration    = newD("Acceleration", "Cursor acceleration value")
	Friction        = newD("Friction", "Cursor friction value")
	WheelAmount     = newD("WheelAmount", "Wheel Up/Down amount(integer)")
	SniperModeSpeed = newD("SniperModeSpeed", "Speed in sniper mode(integer)")
	FlashDistance   = newD("FlashDistance", "Flash distance(integer)")
)
