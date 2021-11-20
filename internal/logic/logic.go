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

var SetCursorPos func(x, y int)
var AddCursorPos func(dx, dy int32)
var GetCursorPos func() (x, y int)
var MouseDown func(button int)
var MouseUp func(button int)
var Wheel func(amount int, hor bool)

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
			lg.Logf("Activate %s", fnc.name)
		}

		if fnc.isActivated && !activate {
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
	AddCursorPos(int32(xSpeed), int32(ySpeed))
	procFriction(&xSpeed)
	procFriction(&ySpeed)
}

func stepFunctions() {
	for _, fnc := range functions {
		if fnc.isActivated {
			fnc.step()
		}
	}
}

func checkKeycode(keycode uint32) (ok bool) {
	_, ok = keycodeState[keycode]
	return
}
