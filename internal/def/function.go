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
	FunctionNameMap[name] = f
	return
}
