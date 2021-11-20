package logic

func GetKeymap() (m map[string][]uint32) {
	functionsMutex.Lock()
	defer functionsMutex.Unlock()
	m = make(map[string][]uint32, len(functions))
	for _, fnc := range functions {
		m[fnc.name] = fnc.keyCodes
	}
	return
}

func SetKeymap(m map[string][]uint32) {
	functionsMutex.Lock()
	defer functionsMutex.Unlock()
	for name, ks := range m {
		for _, fnc := range functions {
			if fnc.name == name {
				fnc.keyCodes = ks
			}
		}
	}
}
