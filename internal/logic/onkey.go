package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/thoas/go-funk"

	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/lg"
)

// todo
func (s *logicState) onKey(keyInfo di.HookKeyInfo) (eat bool) {
	key := keyInfo.Key
	isDown := keyInfo.IsDown
	isRepeat := s.updatePressingKeyMap(key, isDown)
	_ = isRepeat
	s.updateCommandKey(isDown)
	fmt.Printf("%v\n", s.keyState.commandKey)

	return
}

func (s *logicState) updatePressingKeyMap(key string, isDown bool) (isRepeat bool) {
	if isDown {
		isRepeat = funk.ContainsString(s.keyState.pressingKeys, key)
		if !isRepeat {
			s.keyState.pressingKeys = append(s.keyState.pressingKeys, key)
		}
	} else {
		s.keyState.pressingKeys = funk.FilterString(
			s.keyState.pressingKeys, func(s string) bool {
				return s == key
			},
		)
	}
	return
}

func (s *logicState) updateCommandKey(isDown bool) {
	now := time.Now().UnixMilli()
	switch pressingCount := len(s.keyState.pressingKeys); pressingCount {
	case 0: // end last block
		if !isDown {
			s.keyState.lastUpTime = now
		} else {
			lg.Errorf("LOGICAL FALLACY: no pressingKey, but isDown = true")
		}
	case 1:
		if isDown {
			if len(s.keyState.commandKey) > 0 && s.configCache.keyTimeout > now-s.keyState.lastUpTime {

			} else {
				s.keyState.commandKey = [][]string{}
			}
		}
	}
	// replace last block with pressingKeys
}

func printKey(kad di.HookKeyInfo) {
	s := "Up"
	if kad.IsDown {
		s = "Down"
	}
	fmt.Printf("%s %s\n", kad.Key, s)
}

func keyStringToCmdKey(c di.CommandKeyString) (key di.CommandKey) {
	outers := strings.Split(string(c), " - ")
	for _, outer := range outers {
		var inArr []string
		inners := strings.Split(outer, "+")
		for _, inner := range inners {
			inArr = append(inArr, inner)
		}
		key = append(key, inners)
	}
	return
}
