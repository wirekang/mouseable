package typ

type Key string

type OnKeyListener func(key Key, isDown bool) (preventDefault bool)
type OnCursorListener func(x, y int)

type HookManager interface {
	Install()
	Uninstall()
	SetOnKeyListener(listener OnKeyListener)
	SetOnCursorListener(listener OnCursorListener)
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
type DataValue struct {
	int
	float64
	bool
	string
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
	Name() ConfigName
	SetCommandKey(name CommandName, key Key)
	SetDataValue(name DataName, value DataValue)
	CommandKey(name CommandName) Key
	DataValue(name DataName) DataValue
	JSON() ConfigJSON
}

type IOManager interface {
	Save(name ConfigName, data ConfigJSON) error
	Load(name ConfigName) (data ConfigJSON, err error)
	LoadNames() []string
	Lock()
	Unlock()
}

type OverlayManager interface {
	SetVisibility(bool)
	Show()
	Hide()
}

type UIManager interface {
	StartBackground()
	ShowAlert(string)
	ShowError(string)
	SetOnGetNextKeyListener(func() Key)
	SetOnTerminateListener(func())
	SetOnSaveConfigListener(func(ConfigJSON))
	SetOnLoadConfigListener(func(ConfigName) ConfigJSON)
	SetJSONSchema(ConfigJSONSchema)
	SetConfigNames([]ConfigName)
	Open()
}
