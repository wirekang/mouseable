package def

import (
	"encoding/json"

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
	root := map[string]interface{}{}
	root["type"] = "object"
	properties := map[typ.CommandName]interface{}{}
	for _, cmdName := range m.cmdNames {
		properties[cmdName] = map[string]string{
			"type":        "string",
			"description": m.cmdDescriptionMap[cmdName],
		}
	}

	root["properties"] = properties
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
