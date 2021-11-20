package logic

import (
	"fmt"
	"sync"
	"time"

	"github.com/wirekang/vkmap"

	"github.com/wirekang/mouseable/internal/lg"
)

var keycodeState = make(map[uint32]struct{})
var keycodeStateMutex sync.Mutex
var xSpeed, ySpeed int

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
	keycodeStateMutex.Lock()
	defer keycodeStateMutex.Unlock()
	if isDown {
		keycodeState[keyCode] = struct{}{}
	} else {
		delete(keycodeState, keyCode)
	}

	functionsMutex.Lock()
	defer functionsMutex.Unlock()
	for _, fnc := range functions {
		activate := len(fnc.keyCodes) != 0
		for _, kCode := range fnc.keyCodes {
			activate = activate && checkKeycode(kCode)
		}
		if !fnc.isActivated && activate {
			if fnc.onActivate != nil {
				fnc.onActivate()
			}
			lg.Logf("Activate %s", fnc.name)
		}

		if fnc.isActivated && !activate {
			if fnc.onDeactivate != nil {
				fnc.onDeactivate()
			}
			lg.Logf("Deactivate %s", fnc.name)

		}

		fnc.isActivated = activate
		if !preventDefault && activate {
			preventDefault = true
		}
	}

	if lg.IsDev {
		for k := range keycodeState {
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
	for {
		time.Sleep(10 * time.Millisecond)
		moveCursor()
		stepFunctions()
	}
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

func moveCursor() {
	DI.AddCursorPos(int32(xSpeed), int32(ySpeed))
	procFriction(&xSpeed)
	procFriction(&ySpeed)
}

func stepFunctions() {
	functionsMutex.Lock()
	defer functionsMutex.Unlock()
	for _, fnc := range functions {
		if fnc.isActivated && fnc.onStep != nil {
			fnc.onStep()
		}
	}
}

func checkKeycode(keycode uint32) (ok bool) {
	_, ok = keycodeState[keycode]
	return
}
