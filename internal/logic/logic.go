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

// todo
// 고루틴                                                                    참조 데이터
// 1. 주기적으로 계속 커서 이동                                            커서 관련 정보
// 2. 키보드 입력 하면 키 상태 변경 및 명령 실행, 종료 예약, 명령 키바인딩
// isStepping 처리할 필요없음.
// 해당 키가 명령키고 when이 맞는지만 확인하면 무조건 preventDefault
// 3. 주기적으로 실행 예약 확인하면서 명령 begin - step - end 처리
