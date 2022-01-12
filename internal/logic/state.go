package logic

import (
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func (s *logicState) run() {
	s.initListeners()
	s.initConfig()
	s.hookManager.Install()
	go s.keyChanLoop()
	go s.cursorChanLoop()
	go s.cmdLoop()
	go s.cursorLoop()
	s.uiManager.Run()
	defer func() {
		s.ioManager.Unlock()
		lg.Printf("Unlock")
		s.hookManager.Uninstall()
		lg.Printf("Hook uninstalled")
	}()
}

func (s *logicState) initConfig() {
	cn, err := s.ioManager.LoadCurrentConfigName()
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

func (s *logicState) initListeners() {
	s.hookManager.SetKeyInfoChan(s.keyChan, s.preventDefaultChan)
	s.hookManager.SetCursorInfoChan(s.cursorChan)
	s.uiManager.SetOnTerminateListener(s.onTerminate)
	s.uiManager.SetOnGetNextKeyListener(s.onGetNextKey)
	s.uiManager.SetOnSaveConfigListener(s.ioManager.SaveConfig)
	s.uiManager.SetOnLoadConfigListener(s.ioManager.LoadConfig)
	s.uiManager.SetOnLoadConfigSchemaListener(s.definitionManager.JSONSchema)
	s.uiManager.SetOnLoadConfigNamesListener(s.ioManager.LoadConfigNames)
	s.ioManager.SetOnConfigChangeListener(s.onConfigChange)
}

func (s *logicState) onGetNextKey() typ.Key {
	return "Test-Qx2"
}

func (s *logicState) onConfigChange(config typ.Config) {
	lg.Printf("Config changed: %+v", config)
	s.configState.Lock()
	s.configState.doublePressSpeed = int64(config.DataValue("double-press-speed").Int())
	s.overlayManager.SetVisibility(config.DataValue("show-overlay").Bool())
	s.configState.Unlock()
}

func (s *logicState) onTerminate() {
	os.Exit(0)
}

func (s *logicState) cursorChanLoop() {
	for {
		cursorInfo := <-s.cursorChan
		s.overlayManager.SetPosition(cursorInfo.X, cursorInfo.Y)
	}
}
