package logic

import (
	"fmt"
	"math"
	"os"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/di"
	"github.com/wirekang/mouseable/internal/hook"
	"github.com/wirekang/mouseable/internal/io"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/overlay"
	"github.com/wirekang/mouseable/internal/ui"
)

type pointInt struct{ x, y int }
type pointFloat struct{ x, y float64 }

var emptyPointInt = pointInt{}
var emptyPointFloat = pointFloat{}

var emptyStruct = struct{}{}

type logicState struct {
	ioManager         di.IOManager
	hookManager       di.HookManager
	overlayManager    di.OverlayManager
	definitionManager di.DefinitionManager
	uiManager         di.UIManager

	commandTool *di.CommandTool

	cmdState struct {
		when           di.When
		steppingCmdMap map[*di.Command]struct{}
	}

	keyState struct {
		pressingKeyMap map[string]struct{}
		commandKey     di.CommandKey
		lastUpTime     int64
		eatUntilUpMap  map[string]struct{}
		enderMap       map[string]*di.Command
	}

	cursorState struct {
		cursorFixedSpeed    pointInt
		wheelFixedSpeed     pointInt
		cursorSpeed         pointFloat
		wheelSpeed          pointInt
		lastTeleportForward pointInt
	}

	configCache struct {
		keyTimeout         int64
		cursorAcceleration pointFloat
		wheelAcceleration  pointInt
		cursorFriction     pointFloat
		wheelFriction      pointInt
		cursorSniperSpeed  pointInt
		wheelSniperSpeed   pointInt
		teleportDistanceF,
		teleportDistanceX,
		teleportDistanceY int
		commandKeyStringPathMap map[di.CommandKeyString]struct{}
	}

	channel struct {
		configChange chan di.Config

		cursorMove <-chan di.Point

		keyIn  <-chan di.HookKeyInfo
		keyOut chan<- bool

		nextKeyIn  <-chan struct{}
		nextKeyOut chan<- di.CommandKey

		exit <-chan struct{}
	}
	config di.Config
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
			_ = s.configCache.cursorAcceleration.x
			_ = s.configCache.cursorAcceleration.y

			s.cursorState.cursorSpeed.x += 1
			s.cursorState.cursorSpeed.y += 1
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

func Run() {
	uiManager := ui.New()
	defer recoverFn(uiManager)

	ioManager := io.New()
	ioManager.Lock()

	hookManager := hook.New()
	overlayManager := overlay.New()
	definitionManager := def.New()

	logic := logicState{
		ioManager:         ioManager,
		hookManager:       hookManager,
		overlayManager:    overlayManager,
		definitionManager: definitionManager,
		uiManager:         uiManager,
	}

	logic.Run()
}

func recoverFn(uim di.UIManager) {
	cause := recover()
	if cause != nil {
		message := fmt.Sprintf("%v", cause)
		err, ok := cause.(error)
		st := ""
		if ok {
			for {
				stackTracer, ok := err.(interface {
					StackTrace() errors.StackTrace
				})
				if ok {
					st = fmt.Sprintf("%+v\n", stackTracer.StackTrace())
					err = errors.Unwrap(err)
					if err != nil {
						continue
					}
				}

				break
			}
		}
		if st == "" {
			st = string(debug.Stack())
		}

		text := fmt.Sprintf("%s\n\nStackTrace:\n%s", message, st)
		lg.Errorf(text)
		uim.ShowError(text)
		os.Exit(1)
	}
}
