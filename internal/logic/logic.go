package logic

import (
	"fmt"
	"sync"
	"time"

	"github.com/wirekang/vkmap"

	"github.com/wirekang/mouseable/internal/lg"
)

var state struct {
	mutex          sync.Mutex
	keycode        map[uint32]struct{}
	xSpeed, ySpeed int
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
		time.Sleep(10 * time.Millisecond)
		moveCursor()
		stepFunctions()
	}
}

func moveCursor() {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	DI.AddCursorPos(int32(state.xSpeed), int32(state.ySpeed))
	procFriction(&state.xSpeed)
	procFriction(&state.ySpeed)
}

func procFriction(s *int) {
	if *s > 0 {
		*s -= getInt("friction")
		if *s < 0 {
			*s = 0
		}
		return
	}

	if *s < 0 {
		*s += getInt("friction")
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
