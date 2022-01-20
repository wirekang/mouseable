package logic

import (
	"os"
	"os/signal"

	"github.com/wirekang/mouseable/internal/di"
)

func makeConfigChangeListener(c chan<- di.Config) func(di.Config) {
	return func(config di.Config) {
		c <- config
	}
}

func makeCursorListener(c chan<- di.Point) func(di.Point) {
	return func(info di.Point) {
		c <- info
	}
}

func makeKeyListener(c chan<- di.HookKeyInfo, c2 <-chan bool) func(di.HookKeyInfo) bool {
	return func(info di.HookKeyInfo) bool {
		c <- info
		return <-c2
	}
}

func makeOnGetNextKeyListener(
	needNextKeyChan chan<- struct{}, nextKeyChan <-chan di.CommandKeyString,
) func() di.CommandKeyString {
	return func() (key di.CommandKeyString) {
		needNextKeyChan <- emptyStruct
		key = <-nextKeyChan
		return
	}
}

func makeOnExitListener(exitChan chan<- struct{}) func() {
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, os.Interrupt, os.Kill)
		<-sigChan
		exitChan <- emptyStruct
	}()
	return func() {
		exitChan <- emptyStruct
	}
}
