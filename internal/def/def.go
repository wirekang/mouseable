package def

import (
	"encoding/json"
	"fmt"

	"github.com/wirekang/mouseable/internal/typ"
)

type manager struct {
	cmdNames          []typ.CommandName
	cmdOrderMap       map[typ.CommandName]int
	cmdWhenMap        map[typ.CommandName]typ.When
	cmdDescriptionMap map[typ.CommandName]string

	dataNames          []typ.DataName
	dataTypeMap        map[typ.DataName]typ.DataType
	dataDescriptionMap map[typ.DataName]string

	nextFuncOrder int
}

func (m *manager) CommandNames() []typ.CommandName {
	return m.cmdNames
}

func (m *manager) CommandWhen(name typ.CommandName) typ.When {
	return m.cmdWhenMap[name]
}

func (m *manager) DataNames() []typ.DataName {
	return m.dataNames
}

func (m *manager) DataType(name typ.DataName) typ.DataType {
	return m.dataTypeMap[name]
}

func (m *manager) JSONSchema() typ.ConfigJSONSchema {
	command := map[string]interface{}{}
	command["type"] = "object"
	cmdProperties := map[typ.CommandName]interface{}{}
	for _, cmdName := range m.cmdNames {
		cmdProperties[cmdName] = map[string]string{
			"type": "string",
			"description": fmt.Sprintf(
				"%s \n\n when: %s , order: %d",
				m.cmdDescriptionMap[cmdName],
				whenToString(m.cmdWhenMap[cmdName]),
				m.cmdOrderMap[cmdName],
			),
		}
	}
	command["properties"] = cmdProperties

	data := map[string]interface{}{}
	data["type"] = "object"
	dataProperties := map[typ.DataName]interface{}{}
	for _, dataName := range m.dataNames {
		dataProperties[dataName] = map[string]string{
			"type": dataTypeToString(m.dataTypeMap[dataName]),
			"description": fmt.Sprintf(
				"%s \n\n type: %s",
				m.dataDescriptionMap[dataName],
				dataTypeToString(m.dataTypeMap[dataName]),
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
		return typ.ConfigJSONSchema(err.Error())
	}

	return typ.ConfigJSONSchema(s)
}

// nd is new Data
func (m *manager) nd(name, desc string, t typ.DataType) {
	dm := typ.DataName(name)
	m.dataNames = append(m.dataNames, dm)
	m.dataTypeMap[dm] = t
	m.dataDescriptionMap[dm] = desc
}

// nc is new Command
func (m *manager) nc(name, desc string, when typ.When) {
	m.nextFuncOrder += 1
	cn := typ.CommandName(name)
	m.cmdNames = append(m.cmdNames, cn)
	m.cmdOrderMap[cn] = m.nextFuncOrder
	m.cmdDescriptionMap[cn] = desc
	m.cmdWhenMap[cn] = when
}

func whenToString(when typ.When) string {
	switch when {
	case typ.Activated:
		return "Activated"

	case typ.Deactivated:
		return "Deactivated"

	case typ.Any:
		return "Any"
	}

	return "?"
}

func dataTypeToString(dt typ.DataType) string {
	switch dt {
	case typ.Int:
		return "integer"
	case typ.Float:
		return "number"
	case typ.Bool:
		return "boolean"
	case typ.String:
		return "string"
	}
	return "?"
}
