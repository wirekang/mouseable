package logic

import (
	"fmt"
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

func makeCursorListener(c chan<- typ.CursorInfo) func(typ.CursorInfo) {
	return func(info typ.CursorInfo) {
		c <- info
	}
}

func makeKeyListener(c chan<- typ.KeyInfo, c2 <-chan bool) func(typ.KeyInfo) bool {
	return func(info typ.KeyInfo) bool {
		c <- info
		return <-c2
	}
}
func makeOnGetNextKeyListener(needNextKeyChan chan<- struct{}, nextKeyChan <-chan typ.Key) func() typ.Key {
	return func() (key typ.Key) {
		fmt.Println("Start")
		needNextKeyChan <- emptyStruct
		key = <-nextKeyChan
		timoutChan := time.After(time.Second)
		for {
			select {
			case <-timoutChan:
				fmt.Println("Timeout")
				select {
				case <-nextKeyChan:
					fmt.Println("Consume")
				default:
					fmt.Println("Return")
					return
				}

			case needNextKeyChan <- emptyStruct:
				key = <-nextKeyChan
				fmt.Println("Key: ", key)
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
