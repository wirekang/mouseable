package logic

import (
	"os"
	"os/signal"
	"time"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

func makeConfigChangeListener(c chan<- typ.Config) func(typ.Config) {
	return func(config typ.Config) {
		lg.Printf("onConfigChange")
		c <- config
	}
}

func makeCursorListener(c chan<- typ.Point) func(typ.Point) {
	return func(info typ.Point) {
		c <- info
	}
}

func makeKeyListener(c chan<- typ.KeyAndDown, c2 <-chan bool) func(typ.KeyAndDown) bool {
	return func(info typ.KeyAndDown) bool {
		c <- info
		return <-c2
	}
}

func makeOnGetNextKeyListener(needNextKeyChan chan<- struct{}, nextKeyChan <-chan typ.Key) func() typ.Key {
	return func() (key typ.Key) {
		needNextKeyChan <- emptyStruct
		key = <-nextKeyChan
		timoutChan := time.After(time.Second)
		for {
			select {
			case <-timoutChan:
				select {
				case <-nextKeyChan:
				default:
					return
				}

			case needNextKeyChan <- emptyStruct:
				key = <-nextKeyChan
			}
		}
	}
}

func makeOnExitListener(exitChan chan<- struct{}) func() {
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		exitChan <- emptyStruct
	}()
	return func() {
		exitChan <- emptyStruct
	}
}
