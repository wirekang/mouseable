package logic

import (
	"math"
	"sort"
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

func OnUnhook() {
	mutex.Lock()
	defer mutex.Unlock()
	state.steppingMap = map[*logicDef]struct{}{}
}

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
	lg.Logf("OnKey %d %v", keyCode, isDown)
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

	preventDefault = true
	return
}

func Loop() {
	var last int64
	for {
		delta := time.Now().Unix() - last
		last = time.Now().Unix()
		time.Sleep(time.Duration((20 - delta) * 1000))
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
		ix := int32(sx)
		iy := int32(sy)
		if state.fixedSpeed != 0 {
			if ix > 0 {
				ix = int32(state.fixedSpeed)
			} else if ix < 0 {
				ix = -int32(state.fixedSpeed)
			}
			if iy > 0 {
				iy = int32(state.fixedSpeed)
			} else if iy < 0 {
				iy = -int32(state.fixedSpeed)
			}
			state.speedX = float64(ix)
			state.speedY = float64(iy)
		}
		DI.AddCursorPos(ix, iy)
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
	lds := make([]*logicDef, len(state.steppingMap))
	var i int
	for lgc, _ := range state.steppingMap {
		lds[i] = lgc
		i++
	}
	sort.Slice(
		lds, func(i, j int) bool {
			return lds[i].function.Order < lds[i].function.Order
		},
	)
	for _, ld := range lds {
		if ld.onStep != nil {
			ld.onStep(state)
		}
	}
}
