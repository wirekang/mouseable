package logic

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/hook"
	"github.com/wirekang/mouseable/internal/io"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/overlay"
	"github.com/wirekang/mouseable/internal/typ"
	"github.com/wirekang/mouseable/internal/ui"
)

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

	state := state{
		ioManager:          ioManager,
		hookManager:        hookManager,
		overlayManager:     overlayManager,
		definitionManager:  definitionManager,
		uiManager:          uiManager,
		keyChan:            make(chan typ.KeyInfo, 100),
		preventDefaultChan: make(chan bool, 100),
		cursorChan:         make(chan typ.CursorInfo, 100),
		downKeyMap:         make(map[typ.Key]struct{}, 10),
	}

	state.run()
}

func recoverFn(uim typ.UIManager) {
	err := recover()
	if err != nil {
		msg := fmt.Sprintf("panic: %v\n\n", err)
		if st, ok := err.(interface {
			StackTrace() errors.StackTrace
		}); ok {
			msg += fmt.Sprintf("StackTrace: \n%+v", st.StackTrace())
		}
		lg.Errorf(msg)
		uim.ShowError(msg)
		os.Exit(1)
	}
}
