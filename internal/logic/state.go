package logic

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func (s *logicState) Run() {
	s.checkLogics()
	s.init()
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
		case <-s.channels.exit:
			lg.Printf("Exit mainLoop")
			return
		case <-cursorTicker.C:
			s.onCursorTick()
		case point := <-s.channels.cursorMove:
			s.onCursorMove(point.X, point.Y)
		case config := <-s.channels.configChange:
			s.onConfigChange(config)
		case keyAndDown := <-s.channels.keyIn:
			s.channels.keyOut <- s.onKey(keyAndDown)
		case <-cmdStepTicker.C:
			s.onCmdTick()
		}
	}
}

func (s *logicState) onCursorTick() {}

func (s *logicState) onCmdTick() {
	for cmdName := range s.steppingCmdMap {
		cmdLogicMap[cmdName].onStep(s)
	}
	switch s.when {
	case typ.Deactivated:
		s.overlayManager.Hide()
	case typ.Activated:
		s.overlayManager.Show()
	}
}

// todo
func (s *logicState) onKey(keyAndDown typ.KeyAndDown) (preventDefault bool) {
	return
}

func (s *logicState) onConfigChange(config typ.Config) {
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

	for _, configChan := range s.channels.configChanges {
		configChan <- config
	}
}

func (s *logicState) onCursorMove(x, y int) {
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

func (s *logicState) init() {
	s.steppingCmdMap = make(map[typ.CommandName]struct{}, 10)

	keyInfoChan := make(chan typ.KeyAndDown)
	preventDefaultChan := make(chan bool)
	needNextKeyChan := make(chan struct{})
	nextKeyChan := make(chan typ.Key)
	cursorInfoChan := make(chan typ.Point)
	exitChan := make(chan struct{})

	s.channels.keyIn = keyInfoChan
	s.channels.keyOut = preventDefaultChan
	s.channels.nextKeyIn = needNextKeyChan
	s.channels.nextKeyOut = nextKeyChan
	s.channels.cursorMove = cursorInfoChan
	s.channels.exit = exitChan
	s.channels.configChange = make(chan typ.Config)

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

	s.ioManager.SetOnConfigChangeListener(makeConfigChangeListener(s.channels.configChange))
}

func (s *logicState) selectNeedNextKey(combinedKey typ.Key) bool {
	select {
	case <-s.channels.nextKeyIn:
		s.channels.nextKeyOut <- combinedKey
		return true
	default:
		return false
	}
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

func now() int64 {
	return time.Now().UnixMilli()
}
