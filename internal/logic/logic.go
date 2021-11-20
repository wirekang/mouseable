package logic

import (
	"fmt"
	"time"

	"github.com/wirekang/mouseable/internal/lg"
)

var keycodeState = make(map[uint32]struct{})
var xSpeed, ySpeed int

// temp
var friction = 3
var speed = 4

var DI struct {
	SetCursorPos func(x, y int)
	AddCursorPos func(dx, dy int32)
	GetCursorPos func() (x, y int)
	MouseDown    func(button int)
	MouseUp      func(button int)
	Wheel        func(amount int, hor bool)
}

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
	if isDown {
		keycodeState[keyCode] = struct{}{}
	} else {
		delete(keycodeState, keyCode)
	}

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

	fmt.Println(keycodeState)
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
		*s -= friction
		if *s < 0 {
			*s = 0
		}
		return
	}

	if *s < 0 {
		*s += friction
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
