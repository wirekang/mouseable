package logic

import (
	"time"
)

func (s *logicState) cmdLoop() {
	for _ = range time.Tick(time.Millisecond * time.Duration(100)) {
		s.cursorState.Lock()
		s.keyChanState.RLock()
		s.keyChanState.RUnlock()
		s.cursorState.Unlock()
	}
}
