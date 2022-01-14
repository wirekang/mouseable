package logic

import (
	"os"
	"time"

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
	go s.cursorLoop()
	go s.uiManager.Run()
	s.cmdLoop()
	defer func() {
		s.ioManager.Unlock()
		lg.Printf("Unlock")
		s.hookManager.Uninstall()
		lg.Printf("Hook uninstalled")
	}()
}

func (s *logicState) initConfig() {
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
	s.hookManager.SetKeyInfoChan(s.keyChan, s.preventDefaultChan)
	s.hookManager.SetCursorInfoChan(s.cursorChan)
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
	lg.Printf("onGetNextKey")
	timoutChan := time.After(time.Second)
	var key typ.Key
Loop:
	for {
		select {
		case s.needNextKeyChan <- emptyStruct:
		case <-timoutChan:
			lg.Printf("timeout")
			select {
			case <-s.nextKeyChan:
			default:
			}

			break Loop
		case key = <-s.nextKeyChan:
			lg.Printf("next key: %s", key)
			continue
		}
	}

	lg.Printf("final next key: %s", key)

	return key
}

func (s *logicState) onConfigChange(config typ.Config) {
	lg.Printf("Config changed: %+v", config)
	s.overlayManager.SetVisibility(config.DataValue("show-overlay").Bool())

	s.configState.Lock()
	for _, commandName := range s.definitionManager.CommandNames() {
		key := config.CommandKey(commandName)
		if key == "" {
			continue
		}

		s.configState.keyCmdMap[key] = commandName
	}
	s.configState.doublePressSpeed = int64(config.DataValue("double-press-speed").Int())
	s.configState.cursorAccelerationH = config.DataValue("cursor-acceleration-h").Float()
	s.configState.cursorAccelerationV = config.DataValue("cursor-acceleration-v").Float()
	s.configState.cursorFrictionH = config.DataValue("cursor-friction-h").Float()
	s.configState.cursorFrictionV = config.DataValue("cursor-friction-v").Float()
	s.configState.wheelAccelerationH = config.DataValue("wheel-acceleration-h").Int()
	s.configState.wheelAccelerationV = config.DataValue("wheel-acceleration-v").Int()
	s.configState.wheelFrictionH = config.DataValue("wheel-friction-h").Int()
	s.configState.wheelFrictionV = config.DataValue("wheel-friction-v").Int()
	s.configState.sniperModeSpeedH = config.DataValue("sniper-mode-speed-h").Int()
	s.configState.sniperModeSpeedV = config.DataValue("sniper-mode-speed-v").Int()
	s.configState.sniperModeWheelSpeedH = config.DataValue("sniper-mode-wheel-speed-h").Int()
	s.configState.sniperModeWheelSpeedV = config.DataValue("sniper-mode-wheel-speed-v").Int()
	s.configState.teleportDistanceF = config.DataValue("teleport-distance-f").Int()
	s.configState.teleportDistanceH = config.DataValue("teleport-distance-h").Int()
	s.configState.teleportDistanceV = config.DataValue("teleport-distance-v").Int()
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
