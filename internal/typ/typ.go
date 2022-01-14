package typ

type Key string

type KeyInfo struct {
	Key    Key
	IsDown bool
}

type CursorInfo struct {
	X, Y int
}

type HookManager interface {
	Install()
	Uninstall()
	SetKeyInfoChan(chan<- KeyInfo, <-chan bool)
	SetCursorInfoChan(chan<- CursorInfo)
	SetCursorPosition(x, y int)
	AddCursorPosition(dx, dy int)
	CursorPosition() (x, y int)
	MouseDown(button MouseButton)
	MouseUp(button MouseButton)
	Wheel(amount int, isHorizontal bool)
}

type MouseButton uint

const (
	Left   MouseButton = 0
	Middle MouseButton = 1
	Right  MouseButton = 2
)

type CommandName string
type When uint

const (
	Activated   When = 0
	Deactivated When = 1
	Any         When = 2
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
	CommandNames() []CommandName
	CommandWhen(name CommandName) When
	DataNames() []DataName
	DataType(name DataName) DataType
	JSONSchema() ConfigJSONSchema
}

type ConfigName string
type ConfigJSON string
type ConfigJSONSchema string

type Config interface {
	SetCommandKey(name CommandName, key Key)
	SetDataValue(name DataName, value DataValue)
	CommandKey(name CommandName) Key
	DataValue(name DataName) DataValue
	JSON() ConfigJSON
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
	SetOnGetNextKeyListener(func() Key)
	SetOnTerminateListener(func())
	SetOnSaveConfigListener(func(ConfigName, ConfigJSON) error)
	SetOnLoadConfigListener(func(ConfigName) (ConfigJSON, error))
	SetOnLoadConfigSchemaListener(func() ConfigJSONSchema)
	SetOnLoadConfigNamesListener(func() ([]ConfigName, error))
	SetOnLoadAppliedConfigNameListener(func() (ConfigName, error))
	SetOnApplyConfigNameListener(func(name ConfigName) error)
	Open()
}
