package logic

import (
	"fmt"
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
		when di.When
	}

	keyState struct {
		pressingKeys   []string
		commandKey     di.CommandKey
		lastUpTime     int64
		eatUntilUpMap  map[string]struct{}
		steppingCmdMap map[*di.Command]struct{}
		enderMap       map[string]*di.Command
	}

	cursorState struct {
		cursorFixedSpeed pointInt
		wheelFixedSpeed  pointInt
		cursorSpeed      pointFloat
		wheelSpeed       pointInt
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
		},
		Deactivate: func() {
			s.cmdState.when = di.WhenDeactivated
		},
		AccelerateCursor: func(deg float64) {
			// todo
		},
		MouseDown: func(button di.MouseButton) {
			go s.hookManager.MouseDown(button)
		},
		MouseUp: func(button di.MouseButton) {
			go s.hookManager.MouseUp(button)
		},
		MouseWheel: func(isHorizontal bool) {
			// todo
		},
		Teleport: func(deg float64) {
			// todo
		},
		TeleportForward: func() {
			// todo
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
