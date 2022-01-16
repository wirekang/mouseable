package logic

import (
	"time"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func makeConfigChanListener(c chan<- typ.Config) func(typ.Config) {
	return func(config typ.Config) {
		lg.Printf("onConfigChange")
		c <- config
	}
}

func makeCursorChanListener(c chan<- typ.CursorInfo) func(typ.CursorInfo) {
	return func(info typ.CursorInfo) {
		c <- info
	}
}

func makeKeyChanListener(c chan<- typ.KeyInfo, c2 <-chan bool) func(typ.KeyInfo) bool {
	return func(info typ.KeyInfo) bool {
		c <- info
		return <-c2
	}
}
func makeOnGetNextKeyListener(needNextKeyChan chan<- struct{}, nextKeyChan <-chan typ.Key) func() typ.Key {
	return func() (key typ.Key) {
		timoutChan := time.After(time.Second)
		for {
			select {
			case needNextKeyChan <- emptyStruct:
			case <-timoutChan:
				select {
				case <-nextKeyChan:
				default:
					return
				}
			case key = <-nextKeyChan:
			}
		}
	}
}
