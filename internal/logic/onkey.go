package logic

import (
	"time"

	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/lg"
)

func (s *logicState) onKey(keyInfo di.HookKeyInfo) (eat bool) {
	key := keyInfo.Key
	isDown := keyInfo.IsDown
	isRepeat, pressingCount := s.updatePressingKeyMap(key, isDown)
	didBegin := false
	if !isDown || !isRepeat {
		s.updateCommandKey(key, isDown, pressingCount)
		didBegin = s.procCommand(key, isDown)
	}
	eat = s.procEat(key, isDown, didBegin)
	if isDown {
		eat = s.selectNeedNextKey(di.CommandKeyString(key)) || eat
	}
	return
}

func (s *logicState) updatePressingKeyMap(key string, isDown bool) (isRepeat bool, pressingCount int) {
	if isDown {
		_, isRepeat = s.keyState.pressingKeyMap[key]
		if !isRepeat {
			s.keyState.pressingKeyMap[key] = emptyStruct
		}
	} else {
		delete(s.keyState.pressingKeyMap, key)
	}
	pressingCount = len(s.keyState.pressingKeyMap)
	return
}

func (s *logicState) updateCommandKey(key string, isDown bool, pressingCount int) {

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

func (s *logicState) procCommand(key string, isDown bool) (didBegin bool) {
	cmds := s.definitionManager.Command(s.keyState.commandKey, s.cmdState.when)
	for _, cmd := range cmds {
		_, isStepping := s.cmdState.steppingCmdMap[cmd]
		if isDown && !isStepping {
			cmd.OnBegin(s.commandTool)
			s.keyState.enderMap[key] = append(s.keyState.enderMap[key], cmd)
			s.cmdState.steppingCmdMap[cmd] = emptyStruct
			didBegin = true
			lg.Printf("Begin %p by %s", cmd, key)
		}
	}

	if !isDown {
		endCmds := s.keyState.enderMap[key]
		for _, endCmd := range endCmds {
			endCmd.OnEnd(s.commandTool)
			delete(s.keyState.enderMap, key)
			delete(s.cmdState.steppingCmdMap, endCmd)
			lg.Printf("End %p by %s", endCmd, key)
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

func (s *logicState) selectNeedNextKey(cks di.CommandKeyString) bool {
	select {
	case <-s.channel.nextKeyIn:
		s.channel.nextKeyOut <- cks
		return true
	default:
		return false
	}
}
