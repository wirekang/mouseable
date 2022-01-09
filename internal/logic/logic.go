package logic

import (
	"math"
	"sort"
	"sync"
	"time"

	"github.com/thoas/go-funk"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/lg"
)

var mutex = sync.Mutex{}
var functionMap = def.FunctionMap{}
var cachedDataMap = map[*def.DataDefinition]dataCache{}
var lState = logicState{}
var kState = keyState{
	keyCodes: []uint32{},
}

var steppingLogics = []*logicDefinition{}
var willStartLogics = []*logicDefinition{}
var willStopLogics = []*logicDefinition{}
var lastOkMap = map[*logicDefinition]int64{}
var isActivated bool
var wasCursorMoving bool

var winKeys = []uint32{0x5b, 0x5c}
var controlKeys = []uint32{0x11, 0xa2, 0xa3}
var altKeys = []uint32{0x12, 0xa4, 0xa5}
var shiftKeys = []uint32{0x10, 0xa0, 0xa1}
var lastNormalKeyCode uint32

func OnKey(keyCode uint32, isDown bool) (preventDefault bool) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, fd := range def.FunctionDefinitions {
		r := procFunction(keyCode, isDown, fd, functionMap[fd])
		preventDefault = preventDefault || r
	}
	kState = updateKeyState(keyCode, isDown, kState)

	preventDefault = preventDefault || procNormal(keyCode, isDown)
	return
}

var emptyFunctionKey = def.FunctionKey{}

func procFunction(
	keyCode uint32, isDown bool, fd *def.FunctionDefinition, key def.FunctionKey,
) (preventDefault bool) {
	if key == emptyFunctionKey {
		return
	}

	ld := findLogicDefinition(fd)
	if ld == nil {
		lg.Errorf("No LogicDefinition for %s", fd.Name)
		return
	}

	isStepping := funk.Contains(steppingLogics, ld)

	if fd.When == def.Activated && !isActivated || fd.When == def.Deactivated && isActivated {
		if isStepping {
			registerStop(ld)
		}
		return
	}

	wasOk := isKeyOk(kState, key)
	willOk := isKeyOk(updateKeyState(keyCode, isDown, kState), key)

	if !wasOk && !willOk {
		if isStepping {
			lg.Errorf("LOGIC FALLACY: %s is stepping even not ok", ld.function.Name)
			registerStop(ld)
		}
		return
	}

	isDouble := false
	if willOk {
		delta := time.Now().UnixMilli() - lastOkMap[ld]
		if delta <= int64(cachedDataMap[def.DoublePressSpeed].int) {
			isDouble = true
			lastOkMap[ld] = 0
		} else {

			lastOkMap[ld] = time.Now().UnixMilli()
		}
	}

	if key.IsDouble && !isDouble {
		lg.Logf("%s ok but not double", fd.Name)
		return
	}

	if keyCode == key.KeyCode {
		if isDown && willOk || !isDown && wasOk {
			preventDefault = true
		}
	}

	if !wasOk && willOk {
		registerStart(ld)
	}

	if wasOk && !willOk {
		registerStop(ld)
	}

	return
}
func makeModInfo(keyCode uint32) (mi modInfo) {
	mi.isWin = funk.ContainsUInt32(winKeys, keyCode)
	mi.isControl = funk.ContainsUInt32(controlKeys, keyCode)
	mi.isAlt = funk.ContainsUInt32(altKeys, keyCode)
	mi.isShift = funk.ContainsUInt32(shiftKeys, keyCode)
	return
}

func updateModInfo(keyCode uint32, isDown bool, src modInfo) (mi modInfo) {
	newMI := makeModInfo(keyCode)
	mi = src
	if newMI.isWin {
		mi.isWin = isDown
	}
	if newMI.isControl {
		mi.isControl = isDown
	}
	if newMI.isAlt {
		mi.isAlt = isDown
	}
	if newMI.isShift {
		mi.isShift = isDown
	}
	return
}

func registerStart(ld *logicDefinition) {
	if funk.Contains(willStartLogics, ld) {
		lg.Errorf("LOGIC FALLACY: %s is already registered to start", ld.function.Name)
	} else {
		willStartLogics = append(willStartLogics, ld)
	}

}

func registerStop(ld *logicDefinition) {
	if funk.Contains(willStopLogics, ld) {
		lg.Errorf("LOGIC FALLACY: %s is already registered to stop", ld.function.Name)
	} else {
		willStopLogics = append(willStopLogics, ld)
	}

}

func findLogicDefinition(fd *def.FunctionDefinition) *logicDefinition {
	for i := range logicDefinitions {
		if logicDefinitions[i].function == fd {
			return logicDefinitions[i]
		}
	}
	return nil
}

func updateKeyState(keyCode uint32, isDown bool, ks keyState) (newKS keyState) {
	newKS.keyCodes = updateKeyCodes(keyCode, isDown, ks.keyCodes)
	newKS.modInfo = updateModInfo(keyCode, isDown, ks.modInfo)
	return
}

func updateKeyCodes(keyCode uint32, isDown bool, keyCodes []uint32) []uint32 {
	isContain := funk.ContainsUInt32(keyCodes, keyCode)
	if isContain {
		if isDown {
			return keyCodes
		}
		return funk.FilterUInt32(
			keyCodes, func(s uint32) bool {
				return s != keyCode
			},
		)
	} else {
		if isDown {
			return append(keyCodes, keyCode)
		}
		return keyCodes
	}
}

func isKeyOk(ks keyState, key def.FunctionKey) bool {
	return isModOk(key, ks.modInfo) && (funk.ContainsUInt32(ks.keyCodes, key.KeyCode) || key.KeyCode == 0)
}

func isModOk(key def.FunctionKey, mi modInfo) (ok bool) {
	if key.IsWin && !mi.isWin {
		return
	}

	if key.IsControl && !mi.isControl {
		return
	}

	if key.IsAlt && !mi.isAlt {
		return
	}

	if key.IsShift && !mi.isShift {
		return
	}

	ok = true
	return
}

func procNormal(keyCode uint32, isDown bool) (preventDefault bool) {
	mi := makeModInfo(keyCode)
	if isDown && !(mi.isWin || mi.isAlt || mi.isControl || mi.isShift) {
		select {
		case <-DI.NormalKeyChan:
			DI.NormalKeyChan <- keyCode
			lastNormalKeyCode = keyCode
			preventDefault = true
		default:
		}
	}

	if !isDown && lastNormalKeyCode == keyCode {
		lastNormalKeyCode = 0
		preventDefault = true
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
		startFunctions()
		stopFunctions()
		stepFunctions()
		moveCursor()
		rotateWheel()
		procFriction()
		procActivate()
		mutex.Unlock()
	}
}

func startFunctions() {
	for _, ld := range willStartLogics {
		lg.Logf("Start Function %s", ld.function.Name)
		if ld.onStart != nil {
			ld.onStart(&lState)
		}

		if funk.Contains(steppingLogics, ld) {
			lg.Errorf("LOGIC FALLACY: %s already stepping", ld.function.Name)
		} else {
			steppingLogics = append(steppingLogics, ld)
		}
	}
	willStartLogics = []*logicDefinition{}

}

func stopFunctions() {
	for _, ld := range willStopLogics {
		lg.Logf("Stop Function %s", ld.function.Name)
		if ld.onStop != nil {
			ld.onStop(&lState)
		}
		if !funk.Contains(steppingLogics, ld) {
			lg.Errorf("LOGIC FALLACY: %s is not stepping", ld.function.Name)
		}
		steppingLogics = funk.Filter(
			steppingLogics, func(a *logicDefinition) bool {
				return a != ld
			},
		).([]*logicDefinition)
	}
	willStopLogics = []*logicDefinition{}
}

func moveCursor() {
	ix := int(math.Round(lState.cursorDX))
	iy := int(math.Round(lState.cursorDY))
	if ix != 0 || iy != 0 {
		if lState.fixedSpeedH != 0 && lState.fixedSpeedV != 0 {
			if ix > 0 {
				ix = lState.fixedSpeedH
			} else if ix < 0 {
				ix = -lState.fixedSpeedH
			}
			if iy > 0 {
				iy = lState.fixedSpeedV
			} else if iy < 0 {
				iy = -lState.fixedSpeedV
			}
			lState.cursorDX = float64(ix)
			lState.cursorDY = float64(iy)
		}
		DI.AddCursorPos(ix, iy)
		if !wasCursorMoving {
			wasCursorMoving = true
		}
	} else if wasCursorMoving {
		wasCursorMoving = false
	}
}

func rotateWheel() {
	if lState.wheelDX != 0 {
		go DI.Wheel(lState.wheelDX, true)
	}

	if lState.wheelDY != 0 {
		go DI.Wheel(lState.wheelDY, false)
	}
}

func procFriction() {
	procFrictionFloat(cachedDataMap[def.CursorFrictionH].float, &lState.cursorDX)
	procFrictionFloat(cachedDataMap[def.CursorFrictionV].float, &lState.cursorDY)
	procFrictionInt(cachedDataMap[def.WheelFrictionH].int, &lState.wheelDX)
	procFrictionInt(cachedDataMap[def.WheelFrictionV].int, &lState.wheelDY)
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
	cnt := len(steppingLogics)
	if cnt == 0 {
		return
	}

	lds := make([]*logicDefinition, cnt)
	for i := range steppingLogics {
		lds[i] = steppingLogics[i]
	}
	sort.Slice(
		lds, func(i, j int) bool {
			return lds[i].function.Order < lds[i].function.Order
		},
	)
	for _, ld := range lds {
		if ld.onStep != nil {
			ld.onStep(&lState)
		}
	}
}

func procActivate() {
	if lState.willActivate {
		lState.willActivate = false
		isActivated = true
		go DI.OnActivated()
		return
	}

	if lState.willDeactivate {
		lState.willDeactivate = false
		isActivated = false
		steppingLogics = []*logicDefinition{}
		go DI.OnDeactivated()
	}
}
