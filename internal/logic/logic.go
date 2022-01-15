package logic

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/hook"
	"github.com/wirekang/mouseable/internal/io"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/overlay"
	"github.com/wirekang/mouseable/internal/typ"
	"github.com/wirekang/mouseable/internal/ui"
)

var emptyStruct = struct{}{}

type logicState struct {
	ioManager         typ.IOManager
	hookManager       typ.HookManager
	overlayManager    typ.OverlayManager
	definitionManager typ.DefinitionManager
	uiManager         typ.UIManager

	cursorInfoChan <-chan typ.CursorInfo

	configChans []chan<- typ.Config

	cursorSpeedX, cursorSpeedY int
	cursorDX, cursorDY         float64
	wheelSpeedX, wheelSpeedY   int
	wheelDX, wheelDY           int
	willActivate               bool
	willDeactivate             bool

	onConfigChangeChan chan typ.Config

	needNextKeyChan chan<- struct{}
	nextKeyChan     <-chan typ.Key

	internalKeyInfoChan        <-chan typ.KeyInfo
	internalPreventDefaultChan chan<- bool
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
		ioManager:          ioManager,
		hookManager:        hookManager,
		overlayManager:     overlayManager,
		definitionManager:  definitionManager,
		uiManager:          uiManager,
		onConfigChangeChan: make(chan typ.Config),
	}

	logic.Run()
}

func recoverFn(uim typ.UIManager) {
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
