package logic

import (
	"strconv"
)

var data map[string]string

func SetData(keymap map[string][]uint32, d map[string]string) {
	mutex.Lock()
	defer mutex.Unlock()

	data = d
	for name, ks := range keymap {
		for _, fnc := range functions {
			if fnc.name == name {
				fnc.keyCodes = ks
			}
		}
	}
}

func getInt(key string) int {
	i, _ := strconv.ParseInt(data[key], 10, 32)
	return int(i)
}

func getFloat(key string) float64 {
	i, _ := strconv.ParseFloat(data[key], 64)
	return i
}
