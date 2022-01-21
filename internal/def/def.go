package def

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/di"
)

type dataDef struct {
	description string
	dataType    di.DataType
	dflt        interface{}
}

type commandDef struct {
	cmd         *di.Command
	description string
	when        di.When
}

type manager struct {
	keyStringMap map[di.CommandKeyString][]*commandDef
	cmdDefMap    map[di.CommandName]*commandDef
	dataDefMap   map[di.DataName]*dataDef
}

func (m *manager) DataDefault(name di.DataName) di.DataValue {
	data, ok := m.dataDefMap[name]
	defer func() {
		r := recover()
		if r != nil {
			panic(errors.WithStack(fmt.Errorf("%v", r)))
		}
	}()
	if !ok {
		panic(fmt.Sprintf("no data definition for %s", name))
	}

	dv := dataValue{}
	switch data.dataType {
	case di.Int:
		dv.number = float64(data.dflt.(int))
	case di.Float:
		dv.number = data.dflt.(float64)
	case di.Bool:
		dv.bool = data.dflt.(bool)
	case di.String:
		dv.string = data.dflt.(string)
	default:
		panic(fmt.Sprintf("%v is not DataType", data.dataType))
	}
	return dv
}

func (m *manager) SetConfig(config di.Config) {
	m.keyStringMap = make(map[di.CommandKeyString][]*commandDef, len(m.cmdDefMap))
	for commandName := range m.cmdDefMap {
		cks := config.CommandKeyString(commandName)
		m.keyStringMap[cks] = append(m.keyStringMap[cks], m.cmdDefMap[commandName])
	}
}

func (m *manager) Command(key di.CommandKey, when di.When) (r []*di.Command) {
	cks := key.String()
	cmdDefs, ok := m.keyStringMap[cks]
	if !ok {
		return
	}

	for i := range cmdDefs {
		if cmdDefs[i].when != di.WhenAnytime && cmdDefs[i].when != when {
			continue
		}
		r = append(r, cmdDefs[i].cmd)
	}

	return
}

func (m *manager) JSONSchema() di.ConfigJSONSchema {
	command := map[string]interface{}{}
	command["type"] = "object"
	cmdProperties := map[di.CommandName]interface{}{}
	for cmdName, cmdDef := range m.cmdDefMap {
		cmdProperties[cmdName] = map[string]string{
			"type": "string",
			"description": fmt.Sprintf(
				"%s \n\n when: %s",
				cmdDef.description,
				whenToString(cmdDef.when),
			),
			// todo
			// "pattern": "",
		}
	}
	command["properties"] = cmdProperties

	data := map[string]interface{}{}
	data["type"] = "object"
	dataProperties := map[di.DataName]interface{}{}
	for dataName, dataDef := range m.dataDefMap {
		dataProperties[dataName] = map[string]string{
			"type": dataTypeToString(dataDef.dataType),
			"description": fmt.Sprintf(
				"%s \n\n type: %s\ndefault: %v",
				dataDef.description,
				dataTypeToString(dataDef.dataType),
				dataDef.dflt,
			),
		}
	}
	data["properties"] = dataProperties

	root := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": command,
			"data":    data,
		},
	}

	s, err := json.Marshal(root)
	if err != nil {
		return di.ConfigJSONSchema(err.Error())
	}

	return di.ConfigJSONSchema(s)
}

func (m *manager) insertData(name, desc string, t di.DataType, dflt interface{}) {
	dataDef := &dataDef{
		description: desc,
		dataType:    t,
		dflt:        dflt,
	}
	m.dataDefMap[di.DataName(name)] = dataDef
}

func (m *manager) insertCommand(name, description string, when di.When, cmd *di.Command) {
	if cmd.OnBegin == nil {
		cmd.OnBegin = nop
	}
	if cmd.OnStep == nil {
		cmd.OnStep = nop
	}
	if cmd.OnEnd == nil {
		cmd.OnEnd = nop
	}
	m.cmdDefMap[di.CommandName(name)] = &commandDef{
		cmd:         cmd,
		description: description,
		when:        when,
	}
}

func whenToString(when di.When) string {
	switch when {
	case di.WhenActivated:
		return "Activated"

	case di.WhenDeactivated:
		return "Deactivated"

	case di.WhenAnytime:
		return "Any"
	}

	return "?"
}

func dataTypeToString(dt di.DataType) string {
	switch dt {
	case di.Int:
		return "integer"
	case di.Float:
		return "number"
	case di.Bool:
		return "boolean"
	case di.String:
		return "string"
	}
	return "?"
}

func nop(*di.CommandTool) {}

type dataValue struct {
	string string
	bool   bool
	number float64
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
