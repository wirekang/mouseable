package logic

import (
	"strconv"
	"sync"
)

var data map[string]string
var dataMutex sync.Mutex

func SetData(keymap map[string][]uint32, d map[string]string) {
	functionsMutex.Lock()
	defer functionsMutex.Unlock()

	dataMutex.Lock()
	data = d
	dataMutex.Unlock()

	for name, ks := range keymap {
		for _, fnc := range functions {
			if fnc.name == name {
				fnc.keyCodes = ks
			}
		}
	}
}

func getInt(key string) int {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	i, _ := strconv.ParseInt(data[key], 10, 32)
	return int(i)
}

func getFloat(key string) float64 {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	i, _ := strconv.ParseFloat(data[key], 64)
	return i
}
