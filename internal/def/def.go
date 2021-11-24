// Package def declares types for prevent circular dependencies.
package def

type Function struct {
	Order       int
	Name        string
	Description string
}

type Data struct {
	Order       int
	Name        string
	Description string
}

type HotKey struct {
	IsAlt     bool
	IsControl bool
	IsShift   bool
	IsWin     bool
	KeyCode   uint32
}

type Config struct {
	FunctionKeyCodeMap map[*Function]uint32
	DataValueMap       map[*Data]float64
	ActivationKey      HotKey
	DeactivationKey    HotKey
}

var FunctionNameMap = map[string]*Function{}
var DataNameMap = map[string]*Data{}

func newF(order int, name, des string) *Function {
	f := &Function{Order: order, Name: name, Description: des}
	FunctionNameMap[name] = f
	return f
}

func newD(order int, name, des string) *Data {
	d := &Data{Order: order, Name: name, Description: des}
	DataNameMap[name] = d
	return d
}

var (
	MoveRight     = newF(0, "MoveRight", "Move →")
	MoveRightUp   = newF(1, "MoveRightUp", "Move ↗")
	MoveUp        = newF(2, "MoveUp", "Move ↑")
	MoveLeftUp    = newF(3, "MoveLeftUp", "Move ↖")
	MoveLeft      = newF(4, "MoveLeft", "Move ←")
	MoveLeftDown  = newF(5, "MoveLeftDown", "Move ↙")
	MoveDown      = newF(6, "MoveDown", "Move ↓")
	MoveRightDown = newF(7, "MoveRightDown", "Move ↘")
	ClickLeft     = newF(8, "ClickLeft", "Click Left")
	ClickRight    = newF(9, "ClickRight", "Click Right")
	ClickMiddle   = newF(10, "ClickMiddle", "Click Middle")
	WheelUp       = newF(11, "WheelUp", "Wheel Up")
	WheelDown     = newF(12, "WheelDown", "Wheel Down")
	SniperMode    = newF(13, "SniperMode", "Slow down to increase accuracy")
	Flash         = newF(14, "Flash", "Flash cursor to current angle")
)

var (
	Acceleration    = newD(0, "Acceleration", "Cursor acceleration value")
	Friction        = newD(1, "Friction", "Cursor friction value")
	WheelAmount     = newD(2, "WheelAmount", "Wheel Up/Down amount")
	SniperModeSpeed = newD(
		3,
		"SniperModeSpeed", "Speed in sniper mode",
	)
	FlashDistance = newD(4, "FlashDistance", "Flash distance")
)
