package logic

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"

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

	configState  configState
	cursorState  cursorState
	keyChanState keyChanState

	keyChan            chan typ.KeyInfo
	preventDefaultChan chan bool
	cursorChan         chan typ.CursorInfo
	needNextKeyChan    chan struct{}
	nextKeyChan        chan typ.Key
}

type cursorState struct {
	sync.RWMutex
	cursorSpeedX, cursorSpeedY int
	cursorDX, cursorDY         float64
	wheelSpeedX, wheelSpeedY   int
	wheelDX, wheelDY           int
	willActivate               bool
	willDeactivate             bool
}

type keyChanState struct {
	sync.RWMutex
	downKeyMap      map[typ.Key]struct{}
	pressingModKey  typ.Key
	lastDownKey     typ.Key
	lastDownKeyTime int64
}

type configState struct {
	sync.RWMutex
	doublePressSpeed      int64
	keyCmdMap             map[typ.Key]typ.CommandName
	cursorAccelerationH   float64
	cursorAccelerationV   float64
	cursorFrictionH       float64
	cursorFrictionV       float64
	wheelAccelerationH    int
	wheelAccelerationV    int
	wheelFrictionH        int
	wheelFrictionV        int
	sniperModeSpeedH      int
	sniperModeSpeedV      int
	sniperModeWheelSpeedH int
	sniperModeWheelSpeedV int
	teleportDistanceF     int
	teleportDistanceH     int
	teleportDistanceV     int
	showOverlay           bool
}

func Run() {
	uiManager := ui.New()
	defer recoverFn(uiManager)

	ioManager := io.New()
	lg.Printf("IOManager")
	ioManager.Lock()
	lg.Printf("Lock complete")

	hookManager := hook.New()
	lg.Printf("HookManager")
	overlayManager := overlay.New()
	lg.Printf("OverlayManager")
	definitionManager := def.New()
	lg.Printf("DefinitionManager")

	logic := logicState{
		ioManager:          ioManager,
		hookManager:        hookManager,
		overlayManager:     overlayManager,
		definitionManager:  definitionManager,
		uiManager:          uiManager,
		keyChan:            make(chan typ.KeyInfo, 100),
		preventDefaultChan: make(chan bool, 100),
		cursorChan:         make(chan typ.CursorInfo, 100),
		needNextKeyChan:    make(chan struct{}),
		nextKeyChan:        make(chan typ.Key),

		keyChanState: keyChanState{
			downKeyMap: make(map[typ.Key]struct{}, 10),
		},
		configState: configState{
			keyCmdMap: make(map[typ.Key]typ.CommandName, 10),
		},
	}

	logic.run()
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
