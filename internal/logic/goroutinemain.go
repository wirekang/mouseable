package logic

import (
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func (s *logicState) Run() {
	s.initListeners()
	s.hookManager.Install()
	s.loadAndApplyConfig()
	s.initKeyCombineLooper()
	go s.uiManager.Run()
	s.mainLoop()
	defer func() {
		s.ioManager.Unlock()
		lg.Printf("Unlock")
		s.hookManager.Uninstall()
		lg.Printf("Hook uninstalled")
	}()
}

func (s *logicState) mainLoop() {
	cursorTicker := time.NewTicker(time.Millisecond * 33)
	cmdTicker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-cursorTicker.C:
			s.procCursor()
		case cursorInfo := <-s.cursorInfoChan:
			s.onCursor(cursorInfo.X, cursorInfo.Y)
		case config := <-s.onConfigChangeChan:
			s.changeConfig(config)
		case keyInfo := <-s.internalKeyInfoChan:
			s.internalPreventDefaultChan <- s.procKeyInfo(keyInfo)
		case <-cmdTicker.C:
			s.procCmd()
		}
	}
}

func (s *logicState) procCursor() {}

func (s *logicState) procCmd() {}

func (s *logicState) procKeyInfo(keyInfo typ.KeyInfo) (preventDefault bool) {
	return false
}

func (s *logicState) changeConfig(config typ.Config) {
	s.overlayManager.SetVisibility(config.DataValue("show-overlay").Bool())
	for _, configChan := range s.configChans {
		configChan <- config
	}
}

func (s *logicState) onCursor(x, y int) {
	s.overlayManager.SetPosition(x, y)
}

func (s *logicState) initKeyCombineLooper() {
	keyInfoChan := make(chan typ.KeyInfo)
	preventDefaultChan := make(chan bool)
	needNextKeyChan := make(chan struct{})
	nextKeyChan := make(chan typ.Key)
	internalKeyInfoChan := make(chan typ.KeyInfo)
	internalPreventDefaultChan := make(chan bool)
	configChan := make(chan typ.Config)

	keyCombiner := keyCombiner{
		preventKeyUpMap:            make(map[typ.Key]struct{}, 10),
		keyInfoChan:                keyInfoChan,
		preventDefaultChan:         preventDefaultChan,
		needNextKeyChan:            needNextKeyChan,
		nextKeyChan:                nextKeyChan,
		internalKeyInfoChan:        internalKeyInfoChan,
		internalPreventDefaultChan: internalPreventDefaultChan,
		configChan:                 configChan,
	}
	s.configChans = append(s.configChans, configChan)
	s.nextKeyChan = nextKeyChan
	s.needNextKeyChan = needNextKeyChan
	s.hookManager.SetOnKeyListener(makeKeyChanListener(keyInfoChan, preventDefaultChan))
	s.internalKeyInfoChan = internalKeyInfoChan
	s.internalPreventDefaultChan = internalPreventDefaultChan
	go keyCombiner.Run()
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

func (s *logicState) initListeners() {
	cursorInfoChan := make(chan typ.CursorInfo)
	s.cursorInfoChan = cursorInfoChan
	s.hookManager.SetOnCursorMoveListener(makeCursorChanListener(cursorInfoChan))

	s.uiManager.SetOnGetNextKeyListener(makeOnGetNextKeyListener(s.needNextKeyChan, s.nextKeyChan))
	s.uiManager.SetOnTerminateListener(onExit)
	s.uiManager.SetOnSaveConfigListener(s.ioManager.SaveConfig)
	s.uiManager.SetOnLoadConfigListener(s.ioManager.LoadConfig)
	s.uiManager.SetOnLoadConfigSchemaListener(s.definitionManager.JSONSchema)
	s.uiManager.SetOnLoadConfigNamesListener(s.ioManager.LoadConfigNames)
	s.uiManager.SetOnLoadAppliedConfigNameListener(s.ioManager.LoadAppliedConfigName)
	s.uiManager.SetOnApplyConfigNameListener(s.ioManager.ApplyConfig)

	s.ioManager.SetOnConfigChangeListener(makeConfigChanListener(s.onConfigChangeChan))
}

func onExit() {
	os.Exit(0)
}
