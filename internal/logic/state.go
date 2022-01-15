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
	go s.cursorChanLoop()
	s.mainLoop()
	defer func() {
		s.ioManager.Unlock()
		lg.Printf("Unlock")
		s.hookManager.Uninstall()
		lg.Printf("Hook uninstalled")
	}()
}

func (s *logicState) mainLoop() {
	for {
		select {
		case config := <-s.onConfigChangeChan:
			s.changeConfig(config)
			for _, configChan := range s.configChans {
				configChan <- config
			}
		case keyInfo := <-s.internalKeyInfoChan:
			_ = keyInfo
			s.internalPreventDefaultChan <- false
		}
	}
}

func (s *logicState) onConfigChange(config typ.Config) {
	lg.Printf("onConfigChange")
	s.onConfigChangeChan <- config
}

func (s *logicState) changeConfig(config typ.Config) {
	lg.Printf("changeConfig")
	s.overlayManager.SetVisibility(config.DataValue("show-overlay").Bool())
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
	s.hookManager.SetKeyInfoChan(keyInfoChan, preventDefaultChan)
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
	s.hookManager.SetCursorInfoChan(cursorInfoChan)

	s.uiManager.SetOnTerminateListener(s.onTerminate)
	s.uiManager.SetOnGetNextKeyListener(s.onGetNextKey)
	s.uiManager.SetOnSaveConfigListener(s.ioManager.SaveConfig)
	s.uiManager.SetOnLoadConfigListener(s.ioManager.LoadConfig)
	s.uiManager.SetOnLoadConfigSchemaListener(s.definitionManager.JSONSchema)
	s.uiManager.SetOnLoadConfigNamesListener(s.ioManager.LoadConfigNames)
	s.uiManager.SetOnLoadAppliedConfigNameListener(s.ioManager.LoadAppliedConfigName)
	s.uiManager.SetOnApplyConfigNameListener(s.ioManager.ApplyConfig)

	s.ioManager.SetOnConfigChangeListener(s.onConfigChange)
}

func (s *logicState) onGetNextKey() typ.Key {
	timoutChan := time.After(time.Second)
	var key typ.Key
Loop:
	for {
		select {
		case s.needNextKeyChan <- emptyStruct:
		case <-timoutChan:
			select {
			case <-s.nextKeyChan:
			default:
			}

			break Loop
		case key = <-s.nextKeyChan:
			continue
		}
	}
	return key
}

func (s *logicState) onTerminate() {
	os.Exit(0)
}

func (s *logicState) cursorChanLoop() {
	for {
		cursorInfo := <-s.cursorInfoChan
		s.overlayManager.SetPosition(cursorInfo.X, cursorInfo.Y)
	}
}
