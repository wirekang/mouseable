package di

import (
	"strings"
)

type CommandKey [][]string

func (c CommandKey) String() CommandKeyString {
	var outers []string
	for i := range c {
		var inners []string
		for j := range c[i] {
			inners = append(inners, c[i][j])
		}
		outers = append(outers, strings.Join(inners, "+"))
	}
	return CommandKeyString(strings.Join(outers, " - "))
}

type HookKeyInfo struct {
	Key    string
	IsDown bool
}

type Point struct {
	X, Y int
}

type HookManager interface {
	Install()
	Uninstall()
	SetOnKeyListener(func(HookKeyInfo) bool)
	SetOnCursorMoveListener(func(Point))
	SetCursorPosition(x, y int)
	AddCursorPosition(dx, dy int)
	CursorPosition() (x, y int)
	MouseDown(button MouseButton)
	MouseUp(button MouseButton)
	MouseWheel(amount int, isHorizontal bool)
}

type MouseButton uint

const (
	ButtonLeft   MouseButton = 0
	ButtonMiddle MouseButton = 1
	ButtonRight  MouseButton = 2
)

type CommandName string
type When uint

const (
	WhenDeactivated When = 0
	WhenActivated   When = 1
	WhenAnytime     When = 2
)

type DataName string
type DataValue interface {
	Int() int
	Float() float64
	Bool() bool
	String() string
}

type DataType int

const (
	Int    DataType = 0
	Float  DataType = 1
	Bool   DataType = 2
	String DataType = 3
)

type DefinitionManager interface {
	SetConfig(Config)
	Command(CommandKey, When) []*Command
	DataDefault(DataName) DataValue
	JSONSchema() ConfigJSONSchema
}

type ConfigName string
type ConfigJSON string
type ConfigJSONSchema string
type CommandKeyString string

type Config interface {
	CommandKeyString(CommandName) CommandKeyString
	DataValue(DataName) DataValue
	SetJSON(ConfigJSON) error
	JSON() ConfigJSON
	CommandKeyStringPathMap() map[CommandKeyString]struct{}
}

type IOManager interface {
	SaveConfig(ConfigName, ConfigJSON) error
	LoadConfig(ConfigName) (ConfigJSON, error)
	LoadConfigNames() ([]ConfigName, error)
	LoadAppliedConfigName() (ConfigName, error)
	ApplyConfig(ConfigName) error
	Lock()
	Unlock()
	SetOnConfigChangeListener(func(Config))
}

type OverlayManager interface {
	SetVisibility(bool)
	Show()
	Hide()
	SetPosition(x, y int)
}

type UIManager interface {
	Run()
	ShowAlert(string)
	ShowError(string)
	SetOnGetNextKeyListener(func() CommandKeyString)
	SetOnTerminateListener(func())
	SetOnSaveConfigListener(func(ConfigName, ConfigJSON) error)
	SetOnLoadConfigListener(func(ConfigName) (ConfigJSON, error))
	SetOnLoadConfigSchemaListener(func() ConfigJSONSchema)
	SetOnLoadConfigNamesListener(func() ([]ConfigName, error))
	SetOnLoadAppliedConfigNameListener(func() (ConfigName, error))
	SetOnApplyConfigNameListener(func(name ConfigName) error)
	Open()
}

type Direction uint

const (
	DirectionRight     Direction = 0
	DirectionRightUp   Direction = 1
	DirectionUp        Direction = 2
	DirectionLeftUp    Direction = 3
	DirectionLeft      Direction = 4
	DirectionLeftDown  Direction = 5
	DirectionDown      Direction = 6
	DirectionRightDown Direction = 7
)

type CommandTool struct {
	Activate                    func()
	Deactivate                  func()
	RegisterCursorAccelerator   func(direction Direction)
	UnregisterCursorAccelerator func(direction Direction)
	RegisterWheelAccelerator    func(direction Direction)
	UnregisterWheelAccelerator  func(direction Direction)
	FixCursorSpeed              func()
	UnfixCursorSpeed            func()
	FixWheelSpeed               func()
	UnfixWheelSpeed             func()
	MouseDown                   func(button MouseButton)
	MouseUp                     func(button MouseButton)
	Teleport                    func(direction Direction)
	TeleportForward             func()
	Toggle                      func()
}

type Command struct {
	OnBegin func(*CommandTool)
	OnStep  func(*CommandTool)
	OnEnd   func(*CommandTool)
}
