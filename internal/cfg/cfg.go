package cfg

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/typ"
)

type cmdKeyMap map[typ.CommandName]typ.Key
type dataValueMap map[typ.DataName]typ.DataValue

func New(jsn typ.ConfigJSON) (cfg typ.Config, err error) {
	var holder jsonHolder
	err = json.Unmarshal([]byte(jsn), &holder)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	cmdKeyMap := cmdKeyMap{}
	dataValueMap := dataValueMap{}

	for cmd, key := range holder.Command {
		cmdKeyMap[typ.CommandName(cmd)] = typ.Key(key)
	}

	for d, value := range holder.Data {
		dv := dataValue{}
		switch vType := value.(type) {
		case float64:
			dv.isNumber = true
			dv.number = value.(float64)

		case bool:
			dv.isBool = true
			dv.bool = value.(bool)

		case string:
			dv.isString = true
			dv.string = value.(string)

		default:
			err = errors.WithStack(fmt.Errorf("%v is not data type", vType))
			return
		}

		dataValueMap[typ.DataName(d)] = dv
	}

	return &config{
		cmdKeyMap:    cmdKeyMap,
		dataValueMap: dataValueMap,
	}, err
}

type config struct {
	cmdKeyMap    cmdKeyMap
	dataValueMap dataValueMap
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
	r := c.dataValueMap[name]
	if r == nil {
		panic(errors.WithStack(fmt.Errorf("data %s not exitst", name)))
	}
	return r
}

func (c *config) JSON() typ.ConfigJSON {
	cmdMap := map[string]string{}
	dataMap := map[string]interface{}{}
	for cmd, key := range c.cmdKeyMap {
		cmdMap[string(cmd)] = string(key)
	}
	for data, value := range c.dataValueMap {
		dataMap[string(data)] = value
	}

	holder := jsonHolder{
		Command: cmdMap,
		Data:    dataMap,
	}
	b, _ := json.Marshal(holder)
	return typ.ConfigJSON(b)
}
