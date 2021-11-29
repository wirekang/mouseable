package def

type When uint

const (
	Activated   When = 0
	Deactivated When = 1
	Any         When = 2
)

type FunctionDefinition struct {
	Name        string
	Category    string
	Description string
	When        When
	Order       int
}

type FunctionKey struct {
	IsDouble  bool
	IsAlt     bool
	IsControl bool
	IsShift   bool
	IsWin     bool
	KeyCode   uint32
}
