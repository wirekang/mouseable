package logic

import (
	"time"
)

func (s *logicState) cursorLoop() {
	for _ = range time.Tick(time.Millisecond * time.Duration(20)) {
		s.cursorState.Lock()
		s.hookManager.AddCursorPosition(s.procCursorDX(), s.procCursorDY())
		s.hookManager.Wheel(s.procCursorDX(), true)
		s.hookManager.Wheel(s.procCursorDY(), false)
		s.cursorState.Unlock()
	}
}

func (s *logicState) procCursorDX() int {
	return 0
}

func (s *logicState) procCursorDY() int {
	return 0
}

func (s *logicState) procWheelDX() int {
	return 0
}

func (s *logicState) procWheelDY() int {
	return 0
}
