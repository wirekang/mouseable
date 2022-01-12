package logic

import (
	"os"
	"sync"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

type state struct {
	ioManager                  typ.IOManager
	hookManager                typ.HookManager
	overlayManager             typ.OverlayManager
	definitionManager          typ.DefinitionManager
	uiManager                  typ.UIManager
	config                     typ.Config
	cursorSpeedX, cursorSpeedY int
	cursorDX, cursorDY         float64
	wheelSpeedX, wheelSpeedY   int
	wheelDX, wheelDY           int
	willActivate               bool
	willDeactivate             bool
	keyChan                    chan typ.KeyInfo
	preventDefaultChan         chan bool
	cursorChan                 chan typ.CursorInfo
	downKeyMap                 map[typ.Key]int64
	sync.RWMutex
}

func (s *state) run() {
	s.initListeners()
	s.hookManager.Install()
	go s.keyChanLoop()
	go s.cursorChanLoop()
	go s.logicLoop()
	go s.cursorLoop()
	s.uiManager.Run()
	defer func() {
		s.ioManager.Unlock()
		lg.Printf("Unlock")
		s.hookManager.Uninstall()
		lg.Printf("Hook uninstalled")
	}()
}

func (s *state) initListeners() {
	s.hookManager.SetKeyInfoChan(s.keyChan, s.preventDefaultChan)
	s.hookManager.SetCursorInfoChan(s.cursorChan)
	s.uiManager.SetOnTerminateListener(s.onTerminate)
	s.uiManager.SetOnGetNextKeyListener(s.onGetNextKey)
	s.uiManager.SetOnSaveConfigListener(s.ioManager.SaveConfig)
	s.uiManager.SetOnLoadConfigListener(s.ioManager.LoadConfig)
	s.uiManager.SetOnLoadConfigSchemaListener(s.definitionManager.JSONSchema)
	s.uiManager.SetOnLoadConfigNamesListener(s.ioManager.LoadConfigNames)
	s.ioManager.SetOnConfigChangeListener(s.onConfigChange)
}

func (s *state) onGetNextKey() typ.Key {
	return "Test-Qx2"
}

func (s *state) onConfigChange(config typ.Config) {
	s.Lock()
	s.overlayManager.SetVisibility(config.DataValue("ShowOverlay").Bool())
	s.Unlock()
}

func (s *state) onTerminate() {
	os.Exit(0)
}

func (s *state) cursorChanLoop() {
	s.overlayManager.SetVisibility(true)
	s.overlayManager.Show()
	for {
		cursorInfo := <-s.cursorChan
		s.overlayManager.SetPosition(cursorInfo.X, cursorInfo.Y)
	}
}
