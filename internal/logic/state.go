package logic

import (
	"math"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/lg"
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

	stepTicker := time.NewTicker(time.Millisecond * 20)
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
			case <-stepTicker.C:
				s.onStepTick()
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
			s.hookManager.AddCursorPosition(v.x, v.y)
		case v := <-s.channel.wheelBuffer:
			s.hookManager.MouseWheel(v.x, true)
			s.hookManager.MouseWheel(-v.y, false)
		}
	}
}

func (s *logicState) onStepTick() {

	for command := range s.cmdState.steppingCmdMap {
		command.OnStep(s.commandTool)
	}

	if len(s.cursorState.cursorDirectionMap) > 0 {
		s.cursorState.cursorSpeed = combineDirectionMap(
			s.cursorState.cursorDirectionMap, s.cursorState.maxCursorSpeed,
		)
	}

	if len(s.cursorState.wheelDirectionMap) > 0 {
		s.cursorState.wheelSpeed = combineDirectionMap(
			s.cursorState.wheelDirectionMap, s.cursorState.maxWheelSpeed,
		)
	}

	s.channel.cursorBuffer <- s.cursorState.cursorSpeed
	s.channel.wheelBuffer <- s.cursorState.wheelSpeed
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
	s.configCache.cursorSpeed = getInt("cursor-speed")
	s.configCache.wheelSpeed = getInt("wheel-speed")
	s.configCache.cursorSniperSpeed = getInt("cursor-sniper-speed")
	s.configCache.wheelSniperSpeed = getInt("wheel-sniper-speed")
	s.configCache.teleportDistance = getInt("teleport-distance")

	s.configCache.commandKeyStringPathMap = config.CommandKeyStringPathMap()

	s.cursorState.maxCursorSpeed = s.configCache.cursorSpeed
	s.cursorState.maxWheelSpeed = s.configCache.wheelSpeed
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
	s.cursorState.cursorDirectionMap = make(map[di.Direction]struct{}, 8)
	s.cursorState.wheelDirectionMap = make(map[di.Direction]struct{}, 8)
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
	s.channel.cursorBuffer = make(chan vectorInt, 100)
	s.channel.wheelBuffer = make(chan vectorInt, 100)
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

func combineDirectionMap(m map[di.Direction]struct{}, max int) (r vectorInt) {
	var x, y float64
	for direction := range m {
		x += directionVectorMap[direction].x * float64(max)
		y += directionVectorMap[direction].y * float64(max)
	}

	if x == 0 && y == 0 {
		return
	}

	length := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
	x = x / length
	y = y / length
	r.x = int(math.Round(x * float64(max)))
	r.y = int(math.Round(y * float64(max)))
	return
}
