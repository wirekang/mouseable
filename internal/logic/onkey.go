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
	isRepeat, pressingCount := s.updatePressingKeyMap(key, isDown)
	s.updateCommandKey(key, isDown, isRepeat, pressingCount)
	didBegin := s.procCommand(key, isDown, isRepeat)
	eat = s.procEat(key, isDown, didBegin)

	fmt.Printf("%v\n", s.keyState.commandKey)
	return
}

func (s *logicState) updatePressingKeyMap(key string, isDown bool) (isRepeat bool, pressingCount int) {
	if isDown {
		isRepeat = funk.ContainsString(s.keyState.pressingKeys, key)
		if !isRepeat {
			s.keyState.pressingKeys = append(s.keyState.pressingKeys, key)
		}
	} else {
		s.keyState.pressingKeys = funk.FilterString(
			s.keyState.pressingKeys, func(s string) bool {
				return s != key
			},
		)
	}
	pressingCount = len(s.keyState.pressingKeys)
	return
}

func (s *logicState) updateCommandKey(key string, isDown bool, isRepeat bool, pressingCount int) {
	if isDown && isRepeat {
		return
	}

	now := time.Now().UnixMilli()
	if !isDown {
		s.keyState.lastUpTime = now
		return
	}

	didTimeout := s.configCache.keyTimeout < now-s.keyState.lastUpTime
	clone := s.cloneCommandKey()
	switch pressingCount {
	case 0:
		lg.Errorf("LOGICAL FALLACY: isDown=true, pressingCount=0")
		return
	case 1:
		if didTimeout {
			s.keyState.commandKey = [][]string{{key}}
			return
		}

		clone = append(clone, []string{key})
	default:
		if len(clone) < 1 {
			lg.Errorf("LOGICAL FALLACY: isDown=true, pressingCount>1, len(clone) < 1")
			return
		}
		clone[len(clone)-1] = append(clone[len(clone)-1], key)
	}

	if _, ok := s.configCache.commandKeyStringPathMap[clone.String()]; ok {
		s.keyState.commandKey = clone
		return
	} else {
		s.keyState.commandKey = [][]string{{key}}
	}
}

func (s *logicState) procCommand(key string, isDown bool, isRepeat bool) (didBegin bool) {
	if isDown && isRepeat {
		return
	}

	cmd := s.definitionManager.Command(s.keyState.commandKey, s.cmdState.when)
	if cmd != nil {
		_, isStepping := s.keyState.steppingCmdMap[cmd]
		if isDown && !isStepping {
			cmd.OnBegin(s.commandTool)
			s.keyState.enderMap[key] = cmd
			s.keyState.steppingCmdMap[cmd] = emptyStruct
			didBegin = true
			lg.Printf("Begin %d", cmd)
		}
	}

	if !isDown {
		endCmd := s.keyState.enderMap[key]
		if endCmd != nil {
			endCmd.OnEnd(s.commandTool)
			delete(s.keyState.enderMap, key)
			delete(s.keyState.steppingCmdMap, endCmd)
			lg.Printf("End %d", endCmd)
		}
	}
	return
}

func (s *logicState) procEat(key string, isDown bool, didBegin bool) (shouldEat bool) {
	if didBegin {
		s.keyState.eatUntilUpMap[key] = emptyStruct
		shouldEat = true
		return
	}

	_, shouldEat = s.keyState.eatUntilUpMap[key]
	if shouldEat && !isDown {
		delete(s.keyState.eatUntilUpMap, key)
	}
	return
}

func (s *logicState) cloneCommandKey() (clone di.CommandKey) {
	clone = make(di.CommandKey, len(s.keyState.commandKey))
	for i := range s.keyState.commandKey {
		clone[i] = make([]string, len(s.keyState.commandKey[i]))
		copy(clone[i], s.keyState.commandKey[i])
	}
	return
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
