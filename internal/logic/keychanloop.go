package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func (s *logicState) keyChanLoop() {
	for {
		keyInfo := <-s.keyChan
		isDown := keyInfo.IsDown
		originKey := keyInfo.Key
		editedKey, isMod := normalizeModKey(originKey)

		if isMod {
			s.setPressingModKey(editedKey, isDown)
		}

		if isDown {
			if !isMod {
				editedKey = s.combineModKey(editedKey)
			}

			if s.isDouble(originKey) {
				editedKey += "x2"
			}

			cmd, isCmdKey := s.getCmd(editedKey)
			if isCmdKey {
				_ = cmd
			} else {
			}
			s.updateLastDownKey(originKey)
		} else {
		}

		s.updateDownKeyMap(originKey, isDown)

		if isDown {
			fmt.Println(editedKey)
			select {
			case <-s.needNextKeyChan:
				lg.Printf("need next key")
				s.nextKeyChan <- editedKey
				s.preventDefaultChan <- true
			default:
				lg.Printf("not need next key")
				s.preventDefaultChan <- false
			}
		} else {
			s.preventDefaultChan <- false
		}
	}
}

func (s *logicState) setPressingModKey(key typ.Key, isDown bool) {
	s.keyChanState.Lock()
	if isDown {
		s.keyChanState.pressingModKey = key
	} else if s.keyChanState.pressingModKey == key {
		s.keyChanState.pressingModKey = ""
	}
	s.keyChanState.Unlock()
}

// combineModKey convert "A" to "Shift+A" if Shift was pressing.
func (s *logicState) combineModKey(key typ.Key) (rst typ.Key) {
	s.keyChanState.Lock()
	pressingModKey := s.keyChanState.pressingModKey
	s.keyChanState.Unlock()
	if pressingModKey != "" {
		rst = pressingModKey + "+"
	}
	rst += key
	return
}

func (s *logicState) getCmd(key typ.Key) (cmd typ.CommandName, ok bool) {
	s.configState.RLock()
	cmd, ok = s.configState.keyCmdMap[key]
	s.configState.RUnlock()
	return
}

// isDouble converts "A" to "Ax2" if it was pressed recently.
func (s *logicState) isDouble(key typ.Key) (ok bool) {
	s.keyChanState.RLock()
	lastTime := s.keyChanState.lastDownKeyTime
	if s.keyChanState.lastDownKey != key {
		s.keyChanState.RUnlock()
		return
	}

	s.keyChanState.RUnlock()

	s.configState.RLock()
	spd := s.configState.doublePressSpeed
	s.configState.RUnlock()

	ok = time.Now().UnixMilli()-lastTime <= spd
	return
}

func (s *logicState) updateDownKeyMap(key typ.Key, isDown bool) {
	s.keyChanState.Lock()
	if isDown {
		s.keyChanState.downKeyMap[key] = emptyStruct
	} else {
		delete(s.keyChanState.downKeyMap, key)
	}
	s.keyChanState.Unlock()
	return
}

func (s *logicState) updateLastDownKey(key typ.Key) {
	s.keyChanState.Lock()
	s.keyChanState.lastDownKey = key
	s.keyChanState.lastDownKeyTime = time.Now().UnixMilli()
	s.keyChanState.Unlock()
}

// normalizeModKey converts "Left Shift" to "Shift", "Right Ctrl" to "Ctrl". Otherwise, return as it is.
func normalizeModKey(key typ.Key) (v typ.Key, ok bool) {
	v = key
	ok = strings.Contains(string(key), "Alt")
	if ok {
		v = "Alt"
		return
	}

	ok = strings.Contains(string(key), "Shift")
	if ok {
		v = "Shift"
		return
	}

	ok = strings.Contains(string(key), "Ctrl")
	if ok {
		v = "Ctrl"
		return
	}
	return
}
