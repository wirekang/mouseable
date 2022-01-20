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

	cursorTicker := time.NewTicker(time.Millisecond * 20)
	cmdStepTicker := time.NewTicker(time.Millisecond * 200)
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
				s.onCmdStep()
			default:
			}
		}
	}
}

func (s *logicState) onCursorTick() {
	var dx, dy int
	if s.cursorState.cursorFixedSpeed == emptyPointInt {
		dx = int(math.Round(s.cursorState.cursorSpeed.x))
		dy = int(math.Round(s.cursorState.cursorSpeed.y))
	} else {
		dx = s.cursorState.cursorFixedSpeed.x
		dy = s.cursorState.cursorFixedSpeed.y
	}
	go s.hookManager.AddCursorPosition(dx, dy)

	if s.cursorState.wheelFixedSpeed == emptyPointInt {
		dx = s.cursorState.wheelFixedSpeed.x
		dy = s.cursorState.wheelFixedSpeed.y
	} else {
		dx = s.cursorState.wheelSpeed.x
		dy = s.cursorState.wheelSpeed.y
	}
	s.hookManager.MouseWheel(dx, true)
	s.hookManager.MouseWheel(dy, false)
}

func (s *logicState) onCmdStep() {
	for command := range s.cmdState.steppingCmdMap {
		command.OnStep(s.commandTool)
	}

	s.cursorState.cursorSpeed = frictionFloat(s.cursorState.cursorSpeed, s.configCache.cursorFriction)
	s.cursorState.wheelSpeed = frictionInt(s.cursorState.wheelSpeed, s.configCache.wheelFriction)
}

func (s *logicState) onConfigChange(config di.Config) {
	lg.Printf("onConfigChange")
	s.definitionManager.SetConfig(config)

	getInt := s.makeDataGetterInt(config)
	getFloat := s.makeDataGetterFloat(config)
	getBool := s.makeDataGetterBool(config)
	getString := s.makeDataGetterString(config)
	_ = getString

	s.overlayManager.SetVisibility(getBool("show-overlay"))
	s.configCache.keyTimeout = int64(getInt("key-timeout"))
	s.configCache.cursorFriction = pointFloat{
		x: getFloat("cursor-friction-x"),
		y: getFloat("cursor-friction-y"),
	}
	s.configCache.wheelFriction = pointInt{
		x: getInt("wheel-friction-x"),
		y: getInt("wheel-friction-y"),
	}
	s.configCache.cursorAcceleration = pointFloat{
		x: getFloat("cursor-acceleration-x"),
		y: getFloat("cursor-acceleration-y"),
	}
	s.configCache.wheelAcceleration = pointInt{
		x: getInt("wheel-acceleration-x"),
		y: getInt("wheel-acceleration-y"),
	}
	s.configCache.cursorSniperSpeed = pointInt{
		x: getInt("cursor-sniper-speed-x"),
		y: getInt("cursor-sniper-speed-y"),
	}
	s.configCache.wheelSniperSpeed = pointInt{
		x: getInt("wheel-sniper-speed-x"),
		y: getInt("wheel-sniper-speed-y"),
	}
	s.configCache.teleportDistanceF = getInt("teleport-distance-f")
	s.configCache.teleportDistanceX = getInt("teleport-distance-x")
	s.configCache.teleportDistanceY = getInt("teleport-distance-y")

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
	s.cmdState.steppingCmdMap = make(map[*di.Command]struct{}, 5)
	s.keyState.pressingKeyMap = make(map[string]struct{}, 5)
	s.keyState.eatUntilUpMap = make(map[string]struct{}, 5)
	s.keyState.enderMap = make(map[string]*di.Command, 5)

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

func (s *logicState) initCommandTool() {
	s.commandTool = &di.CommandTool{
		Activate: func() {
			s.cmdState.when = di.WhenActivated
			s.overlayManager.Show()
		},
		Deactivate: func() {
			s.cmdState.when = di.WhenDeactivated
			s.overlayManager.Hide()
		},
		AccelerateCursor: func(deg float64) {
			// todo
			s.cursorState.cursorSpeed.x += s.configCache.cursorAcceleration.x
			s.cursorState.cursorSpeed.y += s.configCache.cursorAcceleration.y
		},
		FixCursorSpeed: func() {
			s.cursorState.cursorFixedSpeed = s.configCache.cursorSniperSpeed
		},
		UnfixCursorSpeed: func() {
			s.cursorState.cursorFixedSpeed = emptyPointInt
		},
		FixWheelSpeed: func() {
			s.cursorState.wheelFixedSpeed = s.configCache.wheelSniperSpeed
		},
		UnfixWheelSpeed: func() {
			s.cursorState.wheelFixedSpeed = emptyPointInt
		},
		MouseDown: func(button di.MouseButton) {
			go s.hookManager.MouseDown(button)
		},
		MouseUp: func(button di.MouseButton) {
			go s.hookManager.MouseUp(button)
		},
		MouseWheel: func(deg float64) {},
		Teleport:   func(deg float64) {},
		TeleportForward: func() {
			if math.Abs(s.cursorState.cursorSpeed.x) > 0.3 || math.Abs(s.cursorState.cursorSpeed.y) > 0.3 {
				distance := s.configCache.teleportDistanceF
				angle := math.Atan2(s.cursorState.cursorSpeed.x, s.cursorState.cursorSpeed.y)
				s.cursorState.lastTeleportForward = pointInt{
					x: int(math.Round(float64(distance) * math.Sin(angle))),
					y: int(math.Round(float64(distance) * math.Cos(angle))),
				}
			}
			s.hookManager.AddCursorPosition(
				s.cursorState.lastTeleportForward.x, s.cursorState.lastTeleportForward.y,
			)
		},
	}
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

func frictionInt(p pointInt, f pointInt) (r pointInt) {
	if p.x > 0 {
		r.x = p.x - f.x
		if r.x < 0 {
			r.x = 0
		}
	} else if p.x < 0 {
		r.x = p.x + f.x
		if r.x > 0 {
			r.x = 0
		}
	}

	if p.y > 0 {
		r.y = p.y - f.y
		if r.y < 0 {
			r.y = 0
		}
	} else if p.y < 0 {
		r.y = p.y + f.y
		if r.y > 0 {
			r.y = 0
		}
	}
	return
}

func frictionFloat(p pointFloat, f pointFloat) (r pointFloat) {
	if p.x > 0 {
		r.x = p.x - f.x
		if r.x < 0 {
			r.x = 0
		}
	} else if p.x < 0 {
		r.x = p.x + f.x
		if r.x > 0 {
			r.x = 0
		}
	}

	if p.y > 0 {
		r.y = p.y - f.y
		if r.y < 0 {
			r.y = 0
		}
	} else if p.y < 0 {
		r.y = p.y + f.y
		if r.y > 0 {
			r.y = 0
		}
	}
	return
}
