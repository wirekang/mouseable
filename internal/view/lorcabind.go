package view

import (
	"os/exec"

	"github.com/pkg/errors"
	"github.com/zserge/lorca"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/def"
)

func bindLorca(ui lorca.UI) (err error) {
	a := func(name string, f interface{}) {
		err := ui.Bind("__"+name+"__", f)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	a(
		"loadBind",
		func() interface{} {
			configHolder.Lock()
			defer configHolder.Unlock()
			m := make(map[string]interface{})
			m["functionDefinitions"] = def.FunctionDefinitions
			m["dataDefinitions"] = def.DataDefinitions
			fnm := make(map[string]def.FunctionKey, len(configHolder.FunctionMap))
			for fd := range configHolder.FunctionMap {
				fnm[fd.Name] = configHolder.FunctionMap[fd]
			}
			m["functionNameKeyMap"] = fnm
			dnm := make(map[string]def.DataValue, len(configHolder.DataMap))
			for dd := range configHolder.DataMap {
				dnm[dd.Name] = configHolder.DataMap[dd]
			}
			m["dataNameValueMap"] = dnm
			m["version"] = cnst.VERSION
			return m
		},
	)
	a(
		"getKeyCode",
		func() uint32 {
			DI.NormalKeyChan <- 0
			return <-DI.NormalKeyChan
		},
	)
	a(
		"changeFunction",
		func(name string, key def.FunctionKey) bool {
			configHolder.Lock()
			defer configHolder.Unlock()
			configHolder.FunctionMap[def.FunctionNameMap[name]] = key
			err := DI.SaveConfig(configHolder.Config)
			return err == nil
		},
	)
	a(
		"changeData",
		func(name string, value string) bool {
			configHolder.Lock()
			defer configHolder.Unlock()
			configHolder.DataMap[def.DataNameMap[name]] = def.DataValue(value)
			err := DI.SaveConfig(configHolder.Config)
			return err == nil
		},
	)

	a(
		"openLink",
		func(url string) {
			exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		},
	)

	a(
		"getKeyText",
		func(keyCode uint32) string {
			s, ok := DI.GetKeyText(keyCode)
			if !ok {
				return "<ERROR>"
			}
			return s
		},
	)
	return
}
