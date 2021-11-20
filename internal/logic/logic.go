package logic

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/wirekang/vkmap"

	"github.com/wirekang/mouseable/internal/lg"
)

var state struct {
	mutex          sync.Mutex
	keycode        map[uint32]struct{}
	speedX, speedY float64
	isActivated    bool
}

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	if isDown {
		state.keycode[keyCode] = struct{}{}
	} else {
		delete(state.keycode, keyCode)
	}

	functionsMutex.Lock()
	defer functionsMutex.Unlock()
	for _, fnc := range functions {
		if !state.isActivated && !fnc.isIgnoreDeactivate {
			continue
		}

		isStart := len(fnc.keyCodes) != 0
		for _, kCode := range fnc.keyCodes {
			isStart = isStart && checkKeycode(kCode)
		}
		if !fnc.isStepping && isStart {
			if fnc.onStart != nil {
				fnc.onStart()
			}
			lg.Logf("Start %s", fnc.name)
		}

		if fnc.isStepping && !isStart {
			if fnc.onStop != nil {
				fnc.onStop()
			}
			lg.Logf("Stop %s", fnc.name)
		}

		fnc.isStepping = isStart
		if !preventDefault && isStart {
			preventDefault = true
		}
	}

	if lg.IsDev {
		for k := range state.keycode {
			d := vkmap.Map[k].VK
			if d == "" {
				d = vkmap.Map[k].Description
			}
			fmt.Printf("%d: %s,      ", k, d)
		}
		fmt.Println()
	}
	return
}

func Loop() {
	state.keycode = make(map[uint32]struct{})
	for {
		time.Sleep(15 * time.Millisecond)
		stepFunctions()
		moveCursor()
	}
}

func moveCursor() {
	state.mutex.Lock()
	defer state.mutex.Unlock()

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
	functionsMutex.Lock()
	defer functionsMutex.Unlock()
	for _, fnc := range functions {
		if fnc.isStepping && fnc.onStep != nil {
			fnc.onStep()
		}
	}
}

func checkKeycode(keycode uint32) (ok bool) {
	_, ok = state.keycode[keycode]
	return
}
