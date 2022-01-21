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
	"github.com/wirekang/mouseable/internal/logic/mover"
	"github.com/wirekang/mouseable/internal/overlay"
	"github.com/wirekang/mouseable/internal/ui"
)

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
		enderMap       map[string][]*di.Command
	}

	cursorState struct {
		teleportMover *mover.Mover
		cursorMover   *mover.Mover
		wheelMover    *mover.Mover
	}

	configCache struct {
		keyTimeout              int64
		cursorAccel             float64
		wheelAccel              float64
		cursorMaxSpeed          int
		wheelMaxSpeed           int
		cursorSniperSpeed       int
		wheelSniperSpeed        int
		commandKeyStringPathMap map[di.CommandKeyString]struct{}
	}

	channel struct {
		configChange chan di.Config

		cursorMove <-chan di.Point

		keyIn  <-chan di.HookKeyInfo
		keyOut chan<- bool

		nextKeyIn  <-chan struct{}
		nextKeyOut chan<- di.CommandKeyString

		exit <-chan struct{}

		cursorBuffer chan mover.VectorInt
		wheelBuffer  chan mover.VectorInt
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
