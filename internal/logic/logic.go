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

	keyCmdCacheMap             map[typ.Key]cmdCache
	steppingCmdMap             map[typ.CommandName]struct{}
	cursorSpeedX, cursorSpeedY int
	cursorDX, cursorDY         float64
	wheelSpeedX, wheelSpeedY   int
	wheelDX, wheelDY           int
	when                       typ.When

	onConfigChangeChan chan typ.Config

	cursorInfoChan <-chan typ.CursorInfo

	configChans []chan<- typ.Config

	downedOriginKeyMap   map[typ.Key]struct{}
	downedCombinedKeyMap map[typ.Key]struct{}
	originCombinedKeyMap map[typ.Key]typ.Key
	preventKeyUpMap      map[typ.Key]struct{}
	pressingModKey       typ.Key
	lastDownKey          typ.Key
	lastDownKeyTime      int64

	keyInfoChan        <-chan typ.KeyInfo
	preventDefaultChan chan<- bool

	needNextKeyChan  <-chan struct{}
	nextKeyChan      chan<- typ.Key
	doublePressSpeed int64

	exitChan <-chan struct{}
}

type cmdCache struct {
	name typ.CommandName
	when typ.When
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
		steppingCmdMap:    make(map[typ.CommandName]struct{}, 10),
		when:              typ.Deactivated,

		onConfigChangeChan:   make(chan typ.Config),
		preventKeyUpMap:      make(map[typ.Key]struct{}, 10),
		downedOriginKeyMap:   make(map[typ.Key]struct{}, 10),
		downedCombinedKeyMap: make(map[typ.Key]struct{}, 10),
		originCombinedKeyMap: make(map[typ.Key]typ.Key, 10),
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
