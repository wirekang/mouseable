package logic

import (
	"time"
)

func (s *state) logicLoop() {
	for _ = range time.Tick(time.Millisecond * time.Duration(100)) {
		s.Lock()
		s.Unlock()
	}
}
