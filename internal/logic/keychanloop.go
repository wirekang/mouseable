package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/wirekang/mouseable/internal/typ"
)

func (s *state) keyChanLoop() {
	for {
		keyInfo := <-s.keyChan
		var resultKey string // "Ctrl" "Shift-Qx2" "Q" "W" "E" "Zx2" "Ctrlx2"
		s.Lock()
		key := normalizeModKey(keyInfo.Key)
		didDown := s.updateDownKeyMap(key, keyInfo.IsDown)
		if didDown {
			resultKey = string(key)

		} else {
			s.preventDefaultChan <- s.isStepping(key)
		}

		isDouble := s.isDouble(key)
		if isDouble {
			resultKey += "x2"
		}

		fmt.Println(resultKey)

		s.Unlock()
	}
}

func (s *state) isStepping(key typ.Key) bool { return false }
func (s *state) isDouble(key typ.Key) bool   { return false }

func (s *state) updateDownKeyMap(key typ.Key, isDown bool) (didDown bool) {
	_, wasDown := s.downKeyMap[key]
	if isDown {
		if !wasDown {
			s.downKeyMap[key] = time.Now().UnixMilli()
			didDown = true
		}
	} else {
		delete(s.downKeyMap, key)
	}
	return
}

// normalizeModKey converts "Left Shift" => "Shift", "Right Ctrl" => "Ctrl"
func normalizeModKey(key typ.Key) (v typ.Key) {
	var ok bool
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
