package cfg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/di"
)

func New() di.Config {
	return &config{}
}

type config struct {
	cmdMap  map[di.CommandName]di.CommandKeyString
	dataMap map[di.DataName]di.DataValue
	pathMap map[di.CommandKeyString]struct{}
}

func (c *config) CommandKeyStringPathMap() map[di.CommandKeyString]struct{} {
	return c.pathMap
}

func (c *config) CommandKeyString(name di.CommandName) di.CommandKeyString {
	return di.CommandKeyString(strings.ToLower(string(c.cmdMap[name])))
}

func (c *config) SetJSON(configJSON di.ConfigJSON) (err error) {
	var holder jsonHolder
	err = json.Unmarshal([]byte(configJSON), &holder)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	c.cmdMap = map[di.CommandName]di.CommandKeyString{}
	c.dataMap = map[di.DataName]di.DataValue{}
	c.pathMap = map[di.CommandKeyString]struct{}{}

	for cmd := range holder.Command {
		keyString := holder.Command[cmd]
		c.cmdMap[di.CommandName(cmd)] = di.CommandKeyString(keyString)

		for {
			c.pathMap[di.CommandKeyString(keyString)] = struct{}{}
			lastPlusIndex := strings.LastIndexByte(keyString, '+')
			lastMinusIndex := strings.LastIndexByte(keyString, '-')
			if lastPlusIndex == -1 && lastMinusIndex == -1 {
				break
			}
			var cutIndex int
			if lastMinusIndex > lastPlusIndex {
				cutIndex = lastMinusIndex - 1
			}
			if lastPlusIndex > lastMinusIndex {
				cutIndex = lastPlusIndex
			}

			keyString = keyString[:cutIndex]
		}
	}

	for d, value := range holder.Data {
		dv := dataValue{}
		switch vType := value.(type) {
		case float64:
			dv.number = value.(float64)

		case bool:
			dv.bool = value.(bool)

		case string:
			dv.string = value.(string)

		default:
			err = errors.WithStack(fmt.Errorf("%v is not data type", vType))
			return
		}

		c.dataMap[di.DataName(d)] = dv
	}

	return
}

func (c *config) DataValue(name di.DataName) (v di.DataValue) {
	v = c.dataMap[name]
	return
}

func (c *config) JSON() di.ConfigJSON {
	cmdMap := map[string]string{}
	dataMap := map[string]interface{}{}
	for cmd, key := range c.cmdMap {
		cmdMap[string(cmd)] = string(key)
	}
	for data, value := range c.dataMap {
		dataMap[string(data)] = value
	}

	holder := jsonHolder{
		Command: cmdMap,
		Data:    dataMap,
	}
	b, _ := json.Marshal(holder)
	return di.ConfigJSON(b)
}
