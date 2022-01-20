package logic

import (
	"fmt"
	"math"
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
	defer func() {
		s.ioManager.Unlock()
		lg.Printf("Unlock")
		s.hookManager.Uninstall()
		lg.Printf("Hook uninstalled")
	}()
}

func (s *logicState) mainLoop() {
	s.hookManager.AddCursorPosition(1, 0)

	cursorTicker := time.NewTicker(time.Millisecond * 10)
	cmdStepTicker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case keyInfo := <-s.channel.keyIn:
			s.channel.keyOut <- s.onKey(keyInfo)
		default:
			select {
			case <-s.channel.exit:
				lg.Printf("Exit mainLoop")
				return
			case <-cursorTicker.C:
				s.onCursorTick()
			case point := <-s.channel.cursorMove:
				s.onCursorMove(point.X, point.Y)
			case config := <-s.channel.configChange:
				s.onConfigChange(config)
			case <-cmdStepTicker.C:
				s.onCmdTick()
			default:
			}
		}
	}
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

func (s *logicState) onCursorTick() {
	spd := s.cursorState.cursorSpeed
	if s.cursorState.cursorFixedSpeed != emptyVectorInt {
		spd = minMaxInt(s.cursorState.cursorSpeed, s.cursorState.cursorFixedSpeed)
	}
	s.channel.cursorBuffer <- spd

	spd = s.cursorState.wheelSpeed
	if s.cursorState.wheelFixedSpeed != emptyVectorInt {
		spd = minMaxInt(s.cursorState.wheelSpeed, s.cursorState.wheelFixedSpeed)
	}
	s.channel.wheelBuffer <- spd

}

func (s *logicState) onCmdTick() {
	for command := range s.cmdState.steppingCmdMap {
		command.OnStep(s.commandTool)
	}

	if len(s.cursorState.cursorDirectionMap) > 0 {
		s.cursorState.cursorSpeed = combineDirectionMap(
			s.cursorState.cursorDirectionMap, s.configCache.cursorSpeed,
		)
	}

	if len(s.cursorState.wheelDirectionMap) > 0 {
		s.cursorState.wheelSpeed = combineDirectionMap(
			s.cursorState.wheelDirectionMap, s.configCache.wheelSpeed,
		)
	}
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
	s.configCache.cursorSpeed = vectorInt{
		x: getInt("cursor-speed-x"),
		y: getInt("cursor-speed-y"),
	}
	s.configCache.wheelSpeed = vectorInt{
		x: getInt("wheel-speed-x"),
		y: getInt("wheel-speed-y"),
	}
	s.configCache.cursorSniperSpeed = vectorInt{
		x: getInt("cursor-sniper-speed-x"),
		y: getInt("cursor-sniper-speed-y"),
	}
	s.configCache.wheelSniperSpeed = vectorInt{
		x: getInt("wheel-sniper-speed-x"),
		y: getInt("wheel-sniper-speed-y"),
	}
	s.configCache.teleportDistanceF = getInt("teleport-distance-f")
	s.configCache.teleportDistance = vectorInt{
		x: getInt("teleport-distance-x"),
		y: getInt("teleport-distance-y"),
	}

	s.configCache.commandKeyStringPathMap = config.CommandKeyStringPathMap()
	for keyString := range s.configCache.commandKeyStringPathMap {
		fmt.Println(keyString)
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
	nextKeyChan := make(chan di.CommandKey)
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

func combineDirectionMap(m map[di.Direction]struct{}, v vectorInt) (r vectorInt) {
	var x, y float64
	for direction := range m {
		x += directionVectorMap[direction].x * float64(v.x)
		y += directionVectorMap[direction].y * float64(v.y)
	}
	r = minMaxFloatInt(vectorFloat{x, y}, v)
	if math.Abs(float64(r.x)) == math.Abs(float64(r.y)) {
		r.x = int(math.Round(float64(r.x) * slow))
		r.y = int(math.Round(float64(r.y) * slow))
	}
	return
}

func minMaxFloat(v vectorFloat, mm vectorFloat) (r vectorFloat) {
	r.x = math.Min(math.Max(v.x, -mm.x), mm.x)
	r.y = math.Min(math.Max(v.y, -mm.y), mm.y)
	return
}

func minMaxFloatInt(v vectorFloat, mm vectorInt) (r vectorInt) {
	fr := minMaxFloat(v, vectorFloat{float64(mm.x), float64(mm.y)})
	r.x = int(math.Round(fr.x))
	r.y = int(math.Round(fr.y))
	return
}

func minMaxInt(v vectorInt, mm vectorInt) (r vectorInt) {
	fr := minMaxFloat(vectorFloat{float64(v.x), float64(v.y)}, vectorFloat{float64(mm.x), float64(mm.y)})
	r.x = int(math.Round(fr.x))
	r.y = int(math.Round(fr.y))
	return
}
