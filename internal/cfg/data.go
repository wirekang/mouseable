package cfg

type dataValue struct {
	string                     string
	bool                       bool
	number                     float64
	isNumber, isBool, isString bool
}

func (d dataValue) String() string {
	return d.string
}

func (d dataValue) Bool() bool {
	return d.bool
}

func (d dataValue) Int() int {
	return int(d.number)
}

func (d dataValue) Float() float64 {
	return d.number
}
