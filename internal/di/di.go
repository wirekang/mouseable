package di

type CommandKey [][]string

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
	LeftMouseButton   MouseButton = 0
	MiddleMouseButton MouseButton = 1
	RightMouseButton  MouseButton = 2
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
	Command(CommandKey, When) *Command
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
	SetOnGetNextKeyListener(func() CommandKey)
	SetOnTerminateListener(func())
	SetOnSaveConfigListener(func(ConfigName, ConfigJSON) error)
	SetOnLoadConfigListener(func(ConfigName) (ConfigJSON, error))
	SetOnLoadConfigSchemaListener(func() ConfigJSONSchema)
	SetOnLoadConfigNamesListener(func() ([]ConfigName, error))
	SetOnLoadAppliedConfigNameListener(func() (ConfigName, error))
	Open()
}

type CommandTool struct {
	Activate         func()
	Deactivate       func()
	AccelerateCursor func(deg float64)
	MouseDown        func(button MouseButton)
	MouseUp          func(button MouseButton)
	MouseWheel       func(isHorizontal bool)
	Teleport         func(deg float64)
	TeleportForward  func()
}

type Command struct {
	OnBegin func(*CommandTool)
	OnStep  func(*CommandTool)
	OnEnd   func(*CommandTool)
}
