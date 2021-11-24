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
var dataMap map[*def.DataDef]float64
var state = &logicState{
	steppingMap: map[*logicDef]struct{}{},
}

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
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
		sleep := time.Duration(10-delta) * time.Millisecond
		time.Sleep(sleep)
		mutex.Lock()
		stepFunctions()
		moveCursor()
		procDeactivate()
		mutex.Unlock()
	}
}

func moveCursor() {
	ix := int32(math.Round(state.speedX))
	iy := int32(math.Round(state.speedY))
	if ix != 0 || iy != 0 {
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
		friction := dataMap[def.Friction]
		procFriction(friction, &state.speedX)
		procFriction(friction, &state.speedY)
		if !state.wasCursorMoving {
			state.wasCursorMoving = true
			DI.OnCursorMove()
		}
	} else if state.wasCursorMoving {
		state.wasCursorMoving = false
		DI.OnCursorStop()
	}
}

func procFriction(friction float64, s *float64) {
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
	cnt := len(state.steppingMap)
	if cnt == 0 {
		return
	}

	lds := make([]*logicDef, cnt)
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

func procDeactivate() {
	if state.willDeactivate {
		state.willDeactivate = false
		state.steppingMap = map[*logicDef]struct{}{}
		DI.Unhook()
	}
}
