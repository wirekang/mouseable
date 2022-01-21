package logic

import (
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/logic/mover"
)

func (s *logicState) Run() {
	s.init()
	s.loadAndApplyConfig()
	s.initCommandTool()
	s.hookManager.Install()
	go s.bufferLoop()
	go s.mainLoop()
	s.uiManager.Run()
}

func (s *logicState) mainLoop() {
	s.hookManager.AddCursorPosition(1, 0)

	ticker := time.NewTicker(time.Millisecond * 20)
Loop:
	for {
		select {
		case keyInfo := <-s.channel.keyIn:
			s.channel.keyOut <- s.onKey(keyInfo)
		default:
			select {
			case <-s.channel.exit:
				lg.Printf("Exit mainLoop")
				break Loop
			case point := <-s.channel.cursorMove:
				s.onCursorMove(point.X, point.Y)
			case config := <-s.channel.configChange:
				s.onConfigChange(config)
			case <-ticker.C:
				s.onTick()
			default:
			}
		}
	}

	s.ioManager.Unlock()
	lg.Printf("Unlock")
	s.hookManager.Uninstall()
	lg.Printf("Hook uninstalled")
	os.Exit(0)
}

func (s *logicState) bufferLoop() {
	for {
		select {
		case v := <-s.channel.cursorBuffer:
			s.hookManager.AddCursorPosition(v.X, v.Y)
		case v := <-s.channel.wheelBuffer:
			s.hookManager.MouseWheel(v.X, true)
			s.hookManager.MouseWheel(-v.Y, false)
		}
	}
}

func (s *logicState) onTick() {
	// todo
	// No command use OnStep now.
	// for command := range s.cmdState.steppingCmdMap {
	// 	command.OnStep(s.commandTool)
	// }

	s.cursorState.cursorMover.AddSpeedIfDirection(s.configCache.cursorAccel)
	s.cursorState.wheelMover.AddSpeedIfDirection(s.configCache.wheelAccel)

	s.channel.cursorBuffer <- s.cursorState.cursorMover.Vector()
	s.channel.wheelBuffer <- s.cursorState.wheelMover.Vector()
}

func (s *logicState) onConfigChange(config di.Config) {
	lg.Printf("onConfigChange")
	s.definitionManager.SetConfig(config)

	getInt := s.makeDataGetterInt(config)
	getFloat := s.makeDataGetterFloat(config)
	getBool := s.makeDataGetterBool(config)
	getString := s.makeDataGetterString(config)
	_ = getString
	_ = getFloat

	s.overlayManager.SetVisibility(getBool("show-overlay"))
	s.configCache.keyTimeout = int64(getInt("key-timeout"))
	s.configCache.cursorAccel = getFloat("cursor-acceleration")
	s.configCache.wheelAccel = getFloat("wheel-acceleration")
	s.configCache.cursorMaxSpeed = getInt("cursor-max-speed")
	s.configCache.wheelMaxSpeed = getInt("wheel-max-speed")
	s.configCache.cursorSniperSpeed = getInt("cursor-sniper-speed")
	s.configCache.wheelSniperSpeed = getInt("wheel-sniper-speed")
	s.configCache.commandKeyStringPathMap = config.CommandKeyStringPathMap()

	s.cursorState.cursorMover.SetMaxSpeed(s.configCache.cursorMaxSpeed)
	s.cursorState.wheelMover.SetMaxSpeed(s.configCache.wheelMaxSpeed)

	teleportDistance := getInt("teleport-distance")
	s.cursorState.teleportMover.SetMaxSpeed(teleportDistance)
	s.cursorState.teleportMover.SetSpeed(float64(teleportDistance))

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
	lg.Printf("Apply %s", cn)
}

func (s *logicState) init() {
	s.cursorState.cursorMover = &mover.Mover{}
	s.cursorState.wheelMover = &mover.Mover{}
	s.cursorState.teleportMover = &mover.Mover{}
	s.cmdState.steppingCmdMap = make(map[*di.Command]struct{}, 5)
	s.keyState.pressingKeyMap = make(map[string]struct{}, 5)
	s.keyState.eatUntilUpMap = make(map[string]struct{}, 5)
	s.keyState.enderMap = make(map[string][]*di.Command, 5)

	keyInfoChan := make(chan di.HookKeyInfo)
	eatChan := make(chan bool)
	needNextKeyChan := make(chan struct{})
	nextKeyChan := make(chan di.CommandKeyString)
	cursorInfoChan := make(chan di.Point)
	exitChan := make(chan struct{})

	s.channel.keyIn = keyInfoChan
	s.channel.keyOut = eatChan
	s.channel.nextKeyIn = needNextKeyChan
	s.channel.nextKeyOut = nextKeyChan
	s.channel.cursorMove = cursorInfoChan
	s.channel.exit = exitChan
	s.channel.cursorBuffer = make(chan mover.VectorInt, 100)
	s.channel.wheelBuffer = make(chan mover.VectorInt, 100)
	s.channel.configChange = make(chan di.Config)

	s.hookManager.SetOnCursorMoveListener(makeCursorListener(cursorInfoChan))
	s.hookManager.SetOnKeyListener(makeKeyListener(keyInfoChan, eatChan))

	s.uiManager.SetOnGetNextKeyListener(makeOnGetNextKeyListener(needNextKeyChan, nextKeyChan))
	s.uiManager.SetOnTerminateListener(makeOnExitListener(exitChan))
	s.uiManager.SetOnSaveConfigListener(s.ioManager.SaveConfig)
	s.uiManager.SetOnLoadConfigListener(s.ioManager.LoadConfig)
	s.uiManager.SetOnLoadConfigSchemaListener(s.definitionManager.JSONSchema)
	s.uiManager.SetOnLoadConfigNamesListener(s.ioManager.LoadConfigNames)
	s.uiManager.SetOnLoadAppliedConfigNameListener(s.ioManager.LoadAppliedConfigName)
	s.uiManager.SetOnApplyConfigNameListener(s.ioManager.ApplyConfig)

	s.ioManager.SetOnConfigChangeListener(makeConfigChangeListener(s.channel.configChange))
}

func (s *logicState) dataValueOrDefault(config di.Config, name di.DataName) di.DataValue {
	v := config.DataValue(name)
	if v == nil {
		v = s.definitionManager.DataDefault(name)
		lg.Printf("Use default %s", name)
	}
	return v
}

func (s *logicState) makeDataGetterInt(config di.Config) func(name di.DataName) int {
	return func(name di.DataName) int {
		dv := s.dataValueOrDefault(config, name)
		return dv.Int()
	}
}

func (s *logicState) makeDataGetterFloat(config di.Config) func(name di.DataName) float64 {
	return func(name di.DataName) float64 {
		dv := s.dataValueOrDefault(config, name)
		return dv.Float()
	}
}

func (s *logicState) makeDataGetterBool(config di.Config) func(name di.DataName) bool {
	return func(name di.DataName) bool {
		dv := s.dataValueOrDefault(config, name)
		return dv.Bool()
	}
}

func (s *logicState) makeDataGetterString(config di.Config) func(name di.DataName) string {
	return func(name di.DataName) string {
		dv := s.dataValueOrDefault(config, name)
		return dv.String()
	}
}
