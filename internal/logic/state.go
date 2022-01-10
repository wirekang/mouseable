package logic

import (
	"os"
	"sync"
	"time"

	"github.com/wirekang/mouseable/internal/cnst"
	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

type state struct {
	sync.Mutex
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
}

func (s *state) run() {
	s.initListeners()
	schema := s.definitionManager.JSONSchema()
	s.uiManager.SetJSONSchema(schema)
	s.hookManager.Install()
	go s.loop()
	s.uiManager.StartBackground()
	defer func() {
		s.ioManager.Unlock()
		lg.Printf("Unlock")
		s.hookManager.Uninstall()
		lg.Printf("Hook uninstalled")
	}()
}

func (s *state) initListeners() {
	s.hookManager.SetOnKeyListener(s.onKey)
	s.hookManager.SetOnCursorListener(s.onCursor)
	s.uiManager.SetOnTerminateListener(s.onTerminate)
	s.uiManager.SetOnGetNextKeyListener(s.onGetNextKey)
	s.uiManager.SetOnSaveConfigListener(s.onSave)
	s.uiManager.SetOnLoadConfigListener(s.onLoadConfig)
}

func (s *state) onKey(key typ.Key, isDown bool) (preventDefault bool) {
	logKey(key, isDown)
	return
}

func (s *state) onCursor(x, y int) {
}

func (s *state) onTerminate() {
	os.Exit(0)
}

func (s *state) onGetNextKey() (key typ.Key) {
	return
}

func (s *state) onSave(json typ.ConfigJSON) {}

func (s *state) loop() {
	for _ = range time.Tick(time.Millisecond * time.Duration(20)) {
		s.Lock()
		s.hookManager.AddCursorPosition(s.procCursorDX(), s.procCursorDY())
		s.hookManager.Wheel(s.procCursorDX(), true)
		s.hookManager.Wheel(s.procCursorDY(), false)
		s.Unlock()
	}
}

func (s *state) onLoadConfig(name typ.ConfigName) typ.ConfigJSON {
	data, err := s.ioManager.Load(name)
	if err != nil {
		return typ.ConfigJSON(err.Error())
	}

	return data
}

func logKey(key typ.Key, isDown bool) {
	if !cnst.IsDev {
		return
	}

	t := "Up"
	if isDown {
		t = "Down"
	}
	lg.Printf("'%s' %s\n", key, t)

}
