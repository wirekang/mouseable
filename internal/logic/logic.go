package logic

import (
	"math"
	"sync"
	"time"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/lg"
)

var mutex = sync.Mutex{}
var keyCodeLogicMap = map[uint32]*logicDef{}
var dataMap map[*def.Data]float64
var state = &logicState{
	steppingMap: map[*logicDef]struct{}{},
}

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
	lg.Logf("#### %d %v", keyCode, isDown)
	mutex.Lock()
	defer mutex.Unlock()

	lgc, ok := keyCodeLogicMap[keyCode]
	if !ok {
		return
	}

	_, isStepping := state.steppingMap[lgc]
	isStart := isDown && !isStepping
	isStop := !isDown && isStepping

	if isStart {
		state.steppingMap[lgc] = struct{}{}
		lg.Logf("Start: %s", lgc.function.Name)
		if lgc.onStart != nil {
			lgc.onStart(state)
		}
	}

	if isStop {
		delete(state.steppingMap, lgc)
		lg.Logf("Stop: %s", lgc.function.Name)
		if lgc.onStop != nil {
			lgc.onStop(state)
		}
	}

	preventDefault = isStart || isStop
	lg.Logf("PreventDefault: %d down:%v", keyCode, isDown)

	return
}

func Loop() {
	for {
		time.Sleep(15 * time.Millisecond)
		mutex.Lock()
		stepFunctions()
		moveCursor()
		mutex.Unlock()
	}
}

func moveCursor() {
	sx := math.Round(state.speedX)
	sy := math.Round(state.speedY)
	if sx != 0 || sy != 0 {
		DI.AddCursorPos(int32(sx), int32(sy))
		procFriction(&state.speedX)
		procFriction(&state.speedY)
	}
}

func procFriction(s *float64) {
	friction := dataMap[def.Friction]
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

func stepFunctions() {
	for lgc := range state.steppingMap {
		if lgc.onStep != nil {
			lgc.onStep(state)
		}
	}
}
