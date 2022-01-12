package logic

import (
	"time"
)

func (s *state) cursorLoop() {
	for _ = range time.Tick(time.Millisecond * time.Duration(20)) {
		s.RLock()
		s.hookManager.AddCursorPosition(s.procCursorDX(), s.procCursorDY())
		s.hookManager.Wheel(s.procCursorDX(), true)
		s.hookManager.Wheel(s.procCursorDY(), false)
		s.RUnlock()
	}
}

func (s *state) procCursorDX() int {
	return 0
}

func (s *state) procCursorDY() int {
	return 0
}

func (s *state) procWheelDX() int {
	return 0
}

func (s *state) procWheelDY() int {
	return 0
}
