package cfg

import (
	"github.com/wirekang/mouseable/internal/typ"
)

type cmdKeyMap map[typ.CommandName]typ.Key
type dataValueMap map[typ.DataName]typ.DataValue

func New(name, json string) typ.Config {
	// todo
	return &config{
		name:         typ.ConfigName(name),
		cmdKeyMap:    cmdKeyMap{},
		dataValueMap: dataValueMap{},
	}
}

type config struct {
	name         typ.ConfigName
	cmdKeyMap    cmdKeyMap
	dataValueMap dataValueMap
}

func (c *config) Name() typ.ConfigName {
	return c.name
}

func (c *config) SetCommandKey(name typ.CommandName, key typ.Key) {
	c.cmdKeyMap[name] = key
}

func (c *config) SetDataValue(name typ.DataName, value typ.DataValue) {
	c.dataValueMap[name] = value
}

func (c *config) CommandKey(name typ.CommandName) typ.Key {
	return c.cmdKeyMap[name]
}

func (c *config) DataValue(name typ.DataName) typ.DataValue {
	return c.dataValueMap[name]
}

func (c *config) JSON() typ.ConfigJSON {
	// todo
	return "{}"
}
