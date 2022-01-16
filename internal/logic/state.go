package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func (s *logicState) Run() {
	s.checkLogics()
	s.initListenersAndChannels()
	s.hookManager.Install()
	s.loadAndApplyConfig()
	go s.uiManager.Run()
	s.mainLoop()
	defer func() {
		s.ioManager.Unlock()
		lg.Printf("Unlock")
		s.hookManager.Uninstall()
		lg.Printf("Hook uninstalled")
	}()
}

func (s *logicState) checkLogics() {
	cmdNames := s.definitionManager.CommandNames()
Loop:
	for nameInLogic := range cmdLogicMap {
		for _, nameInNames := range cmdNames {
			if nameInNames == nameInLogic {
				continue Loop
			}
		}
		lg.Errorf("No definition for logic %s", nameInLogic)
	}

	for _, cmdName := range cmdNames {
		lgc, ok := cmdLogicMap[cmdName]
		if !ok {
			lg.Errorf("No logic for definition %s", cmdName)
			cmdLogicMap[cmdName] = cmdLogic{
				onBegin: nop,
				onStep:  nop,
				onEnd:   nop,
			}
		} else {
			if lgc.onBegin == nil {
				lgc.onBegin = nop
			}
			if lgc.onStep == nil {
				lgc.onStep = nop
			}
			if lgc.onEnd == nil {
				lgc.onEnd = nop
			}
			cmdLogicMap[cmdName] = lgc
		}
	}
}

func (s *logicState) mainLoop() {
	cursorTicker := time.NewTicker(time.Millisecond * 33)
	cmdStepTicker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-s.exitChan:
			lg.Printf("Exit mainLoop")
			return
		case <-cursorTicker.C:
			s.procCursor()
		case cursorInfo := <-s.cursorInfoChan:
			s.onCursor(cursorInfo.X, cursorInfo.Y)
		case config := <-s.onConfigChangeChan:
			s.changeConfig(config)
		case keyInfo := <-s.keyInfoChan:
			s.preventDefaultChan <- s.procKeyInfo(keyInfo)
		case <-cmdStepTicker.C:
			s.procCmdStep()
		}
	}
}

func (s *logicState) procCursor() {}

func (s *logicState) procCmdStep() {
	for cmdName := range s.steppingCmdMap {
		cmdLogicMap[cmdName].onStep(s)
	}
	switch s.when {
	case typ.Activated:
		s.overlayManager.Show()
	case typ.Deactivated:
		s.overlayManager.Hide()
	}
}

func (s *logicState) procKeyInfo(keyInfo typ.KeyInfo) (preventDefault bool) {
	originKey := keyInfo.Key
	isDown := keyInfo.IsDown
	_, wasOriginDown := s.downedOriginKeyMap[originKey]
	isFirstDown := isDown && !wasOriginDown

	normalizedKey, isMod := normalizeModKey(originKey)

	// update pressingModKey
	if isMod {
		if isDown {
			s.pressingModKey = normalizedKey
		} else {
			if s.pressingModKey == normalizedKey {
				s.pressingModKey = ""
			}
		}
	}

	combinedKey := normalizedKey

	// combine with mod key
	if !isMod {
		if isDown {
			if s.pressingModKey != "" {
				combinedKey = s.pressingModKey + "+" + normalizedKey
			}
		}
	}

	// check is double press
	if isFirstDown {
		if s.lastDownKey == normalizedKey {
			if now()-s.lastDownKeyTime < s.doublePressSpeed {
				combinedKey += "x2"
			}
		}
		s.lastDownKey = normalizedKey
		s.lastDownKeyTime = now()
	}

	// send next key to front if needed.
	if isDown && s.procNextKey(combinedKey) {
		preventDefault = true
	}

	_, wasCombinedDown := s.downedCombinedKeyMap[combinedKey]

	// key up, but was not downed, is combined key.
	if !isDown && !wasCombinedDown {
		if ck, ok := s.popOriginCombinedMap(originKey); ok {
			combinedKey = ck
		}
	}

	if isDown {
		s.originCombinedKeyMap[originKey] = combinedKey
		s.downedCombinedKeyMap[combinedKey] = emptyStruct
	} else {
		delete(s.downedCombinedKeyMap, combinedKey)
	}

	dt := "Up"
	if isDown {
		dt = "Down"
	}
	fmt.Printf("%s %s\n", combinedKey, dt)

	cmdCache, ok := s.keyCmdCacheMap[combinedKey]
	if ok {
		isSameWhen := s.when == cmdCache.when
		if isSameWhen {
			preventDefault = true
		}

		_, isStepping := s.steppingCmdMap[cmdCache.name]

		if isSameWhen && keyInfo.IsDown && !isStepping {
			lg.Printf("Begin %s", cmdCache.name)
			cmdLogicMap[cmdCache.name].onBegin(s)
			s.steppingCmdMap[cmdCache.name] = emptyStruct
		}

		if !keyInfo.IsDown && isStepping {
			lg.Printf("End %s", cmdCache.name)
			cmdLogicMap[cmdCache.name].onEnd(s)
			delete(s.steppingCmdMap, cmdCache.name)
		}
	}

	// update downedOriginKeyMap
	if isDown {
		s.downedOriginKeyMap[originKey] = emptyStruct
	} else {
		delete(s.downedOriginKeyMap, originKey)
	}

	if isDown && preventDefault {
		s.preventKeyUpMap[originKey] = emptyStruct
	}
	if !isDown && s.popPreventKeyUpMap(originKey) {
		preventDefault = true
	}

	return
}

func (s *logicState) changeConfig(config typ.Config) {
	s.overlayManager.SetVisibility(config.DataValue("show-overlay").Bool())
	s.doublePressSpeed = int64(config.DataValue("double-press-speed").Int())
	s.keyCmdCacheMap = make(map[typ.Key]cmdCache, 20)
	for _, cmdName := range s.definitionManager.CommandNames() {
		key := config.CommandKey(cmdName)
		if key == "" {
			continue
		}

		s.keyCmdCacheMap[key] = cmdCache{
			name: cmdName,
			when: s.definitionManager.CommandWhen(cmdName),
		}
	}

	for _, configChan := range s.configChans {
		configChan <- config
	}
}

func (s *logicState) onCursor(x, y int) {
	s.overlayManager.SetPosition(x, y)
}

func (s *logicState) loadAndApplyConfig() {
	cn, err := s.ioManager.LoadAppliedConfigName()
	if err != nil {
		err = errors.WithStack(err)
		panic(err)
	}

	err = s.ioManager.ApplyConfig(cn)
	if err != nil {
		err = errors.WithStack(err)
		panic(err)
	}
}

func (s *logicState) initListenersAndChannels() {
	keyInfoChan := make(chan typ.KeyInfo)
	preventDefaultChan := make(chan bool)
	needNextKeyChan := make(chan struct{})
	nextKeyChan := make(chan typ.Key)
	cursorInfoChan := make(chan typ.CursorInfo)
	exitChan := make(chan struct{})

	s.keyInfoChan = keyInfoChan
	s.preventDefaultChan = preventDefaultChan
	s.needNextKeyChan = needNextKeyChan
	s.nextKeyChan = nextKeyChan
	s.cursorInfoChan = cursorInfoChan
	s.exitChan = exitChan

	s.hookManager.SetOnCursorMoveListener(makeCursorListener(cursorInfoChan))
	s.hookManager.SetOnKeyListener(makeKeyListener(keyInfoChan, preventDefaultChan))

	s.uiManager.SetOnGetNextKeyListener(makeOnGetNextKeyListener(needNextKeyChan, nextKeyChan))
	s.uiManager.SetOnTerminateListener(makeOnExitListener(exitChan))
	s.uiManager.SetOnSaveConfigListener(s.ioManager.SaveConfig)
	s.uiManager.SetOnLoadConfigListener(s.ioManager.LoadConfig)
	s.uiManager.SetOnLoadConfigSchemaListener(s.definitionManager.JSONSchema)
	s.uiManager.SetOnLoadConfigNamesListener(s.ioManager.LoadConfigNames)
	s.uiManager.SetOnLoadAppliedConfigNameListener(s.ioManager.LoadAppliedConfigName)
	s.uiManager.SetOnApplyConfigNameListener(s.ioManager.ApplyConfig)

	s.ioManager.SetOnConfigChangeListener(makeConfigChangeListener(s.onConfigChangeChan))
}

func normalizeModKey(key typ.Key) (v typ.Key, isModKey bool) {
	v = key
	isModKey = strings.Contains(string(key), "Alt")
	if isModKey {
		v = "Alt"
		return
	}

	isModKey = strings.Contains(string(key), "Shift")
	if isModKey {
		v = "Shift"
		return
	}

	isModKey = strings.Contains(string(key), "Ctrl")
	if isModKey {
		v = "Ctrl"
		return
	}
	return
}

func (s *logicState) procNextKey(combinedKey typ.Key) bool {
	select {
	case <-s.needNextKeyChan:
		s.nextKeyChan <- combinedKey
		return true
	default:
		return false
	}
}

func (s *logicState) popPreventKeyUpMap(key typ.Key) (ok bool) {
	_, ok = s.preventKeyUpMap[key]
	if ok {
		delete(s.preventKeyUpMap, key)
	}
	return
}

func (s *logicState) popOriginCombinedMap(origin typ.Key) (combined typ.Key, ok bool) {
	combined, ok = s.originCombinedKeyMap[origin]
	if ok {
		delete(s.preventKeyUpMap, origin)
	}
	return
}

func now() int64 {
	return time.Now().UnixMilli()
}
