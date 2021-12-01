package logic

import (
	"math"
	"sort"
	"sync"
	"time"

	"github.com/wirekang/mouseable/internal/def"
)

var mutex = sync.Mutex{}
var functionMap = def.FunctionMap{}
var cachedDataMap = map[*def.DataDefinition]dataCache{}
var state = &logicState{
	steppingLogics: []*logicDefinition{},
}

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
	mutex.Lock()
	defer mutex.Unlock()

	if isDown && true { // todo
		select {
		case <-DI.NormalKeyChan:
			DI.NormalKeyChan <- keyCode
		default:
		}
	}

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
		rotateWheel()
		procFriction()
		procDeactivate()
		mutex.Unlock()
	}
}

func moveCursor() {
	ix := int32(math.Round(state.cursorDX))
	iy := int32(math.Round(state.cursorDY))
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
			state.cursorDX = float64(ix)
			state.cursorDY = float64(iy)
		}
		DI.AddCursorPos(ix, iy)
		if !state.wasCursorMoving {
			state.wasCursorMoving = true
			DI.OnCursorMove()
		}
	} else if state.wasCursorMoving {
		state.wasCursorMoving = false
		DI.OnCursorStop()
	}
}

func rotateWheel() {
	if state.wheelDX != 0 {
		DI.Wheel(state.wheelDY, true)
	}

	if state.wheelDY != 0 {
		DI.Wheel(state.wheelDY, false)
	}
}

func procFriction() {
	cursorFriction := cachedDataMap[def.CursorFriction].float
	procFrictionFloat(cursorFriction, &state.cursorDX)
	procFrictionFloat(cursorFriction, &state.cursorDY)
	wheelFriction := cachedDataMap[def.WheelFriction].int
	procFrictionInt(wheelFriction, &state.wheelDX)
	procFrictionInt(wheelFriction, &state.wheelDY)
}

func procFrictionFloat(friction float64, s *float64) {
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

func procFrictionInt(friction int, s *int) {
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
	cnt := len(state.steppingLogics)
	if cnt == 0 {
		return
	}

	lds := make([]*logicDefinition, cnt)
	for i := range state.steppingLogics {
		lds[i] = state.steppingLogics[i]
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
		state.steppingLogics = []*logicDefinition{}
		DI.OnDeactivated()
	}
}
