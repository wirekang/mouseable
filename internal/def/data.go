package def

import (
	"strconv"
	"strings"
)

type DataDefinition struct {
	Name        string
	Description string
	Type        Type
}

type Type int

const (
	Int    Type = 0
	Float  Type = 1
	Bool   Type = 2
	String Type = 3
)

type DataValue string

func (d DataValue) Int() (r int) {
	i, _ := strconv.ParseInt(string(d), 10, 32)
	r = int(i)
	return
}

func (d DataValue) Float() (r float64) {
	r, _ = strconv.ParseFloat(string(d), 64)
	return
}

func (d DataValue) Bool() (r bool) {
	r = strings.ToLower(strings.TrimSpace(string(d))) == "true"
	return
}

func nD(name, desc string, t Type) (d *DataDefinition) {
	d = new(DataDefinition)
	d.Name = name
	d.Description = desc
	d.Type = t
	DataDefinitions = append(DataDefinitions, d)
	DataNameMap[name] = d
	return
}
