package logic

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/wirekang/vkmap"

	"github.com/wirekang/mouseable/internal/lg"
)

var mutex = sync.Mutex{}

var state struct {
	keyCodeMap     map[uint32]struct{}
	speedX, speedY float64
	isActivated    bool
}

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
	lg.Logf("#### %d %v", keyCode, isDown)
	mutex.Lock()
	defer mutex.Unlock()

	if lg.IsDev {
		lg.Logf("BEFORE")
		logState()
	}

	setMapKey(state.keyCodeMap, keyCode, isDown)

	if lg.IsDev {
		lg.Logf("AFTER")
		logState()
	}

	for _, fnc := range functions {
		if !state.isActivated && !fnc.isIgnoreDeactivate {
			if fnc.isStepping {
				fnc.isStepping = false
				if fnc.onStop != nil {
					fnc.onStop()
				}
			}
			continue
		}

		isAllKeysInState := isContainsAll(fnc.keyCodes, state.keyCodeMap)
		isStart := isDown && !fnc.isStepping && isAllKeysInState
		isStop := !isDown && fnc.isStepping && !isAllKeysInState

		if !preventDefault &&
			(((isStart || isStop) && (isDown == isStart) && (isDown == !isStop)) ||
				isAllKeysInState && isContains(
					keyCode, fnc.keyCodes,
				)) &&
			!(state.isActivated && fnc.isIgnoreDeactivate && !isDown) {
			lg.Logf(
				"Prevent %d isDown: %v isStart: %v isStop: %v", keyCode, isDown,
				isStart, isStop,
			)
			preventDefault = true
		}

		if isStart {
			if fnc.onStart != nil {
				fnc.onStart()
			}
			lg.Logf("Start %s", fnc.name)
		}

		if isStop {
			if fnc.onStop != nil {
				fnc.onStop()
			}
			lg.Logf("Stop %s", fnc.name)
		}

		fnc.isStepping = isAllKeysInState
	}

	return
}

func Loop() {
	initState()
	for {
		time.Sleep(15 * time.Millisecond)
		mutex.Lock()
		stepFunctions()
		moveCursor()
		mutex.Unlock()
	}
}

func initState() {
	state.keyCodeMap = map[uint32]struct{}{}
}

func moveCursor() {

	DI.AddCursorPos(
		int32(math.Round(state.speedX)), int32(math.Round(state.speedY)),
	)

	procFriction(&state.speedX)
	procFriction(&state.speedY)
}

func procFriction(s *float64) {
	if *s > 0 {
		*s -= getFloat("friction")
		if *s < 0 {
			*s = 0
		}
		return
	}

	if *s < 0 {
		*s += getFloat("friction")
		if *s > 0 {
			*s = 0
		}
	}
}

func stepFunctions() {
	for _, fnc := range functions {
		if fnc.isStepping {
			if fnc.onStep != nil {
				fnc.onStep()
			}
		}
	}
}

func isContainsAll(slice []uint32, m map[uint32]struct{}) (rst bool) {
	rst = len(slice) != 0
	for _, ui := range slice {
		_, ok := m[ui]
		rst = rst && ok
	}
	return
}

func isContains(ui uint32, slice []uint32) (ok bool) {
	for i := range slice {
		if ui == slice[i] {
			ok = true
			return
		}
	}
	return
}

func setMapKey(m map[uint32]struct{}, key uint32, isSet bool) {
	if isSet {
		m[key] = struct{}{}
	} else {
		delete(m, key)
	}
}

func logState() {
	for k := range state.keyCodeMap {
		d := vkmap.Map[k].VK
		if d == "" {
			d = vkmap.Map[k].Description
		}
		fmt.Printf("%d: %s,      ", k, d)
	}
	fmt.Println()
}
