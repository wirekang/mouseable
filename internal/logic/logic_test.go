package logic_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/logic"
)

type xy struct {
	x, y int
}

var setCursorPosChan = make(chan xy, 100)
var addCursorPosChan = make(chan xy, 100)
var getCursorPosChan = make(chan xy, 100)
var mouseDownChan = make(chan int, 100)
var mouseUpChan = make(chan int, 100)
var mouseWheelChan = make(chan int, 100)
var onActivatedChan = make(chan struct{}, 100)
var onDeactivatedChan = make(chan struct{}, 100)
var normalKeyChan = make(chan uint32)

var control uint32 = 0xa2
var alt uint32 = 0xa4
var shift uint32 = 0xa0
var win uint32 = 0x5b

func init() {
	logic.DI = struct {
		SetCursorPos  func(x int, y int)
		AddCursorPos  func(dx int, dy int)
		GetCursorPos  func() (x int, y int)
		MouseDown     func(button int)
		MouseUp       func(button int)
		Wheel         func(amount int, hor bool)
		OnCursorMove  func()
		OnCursorStop  func()
		OnActivated   func()
		OnDeactivated func()
		NormalKeyChan chan uint32
	}{
		SetCursorPos: func(x int, y int) {
			setCursorPosChan <- xy{x, y}
		}, AddCursorPos: func(dx int, dy int) {
			addCursorPosChan <- xy{dx, dy}
		}, GetCursorPos: func() (x int, y int) {
			xy := <-getCursorPosChan
			x = xy.x
			y = xy.y
			return
		}, MouseDown: func(button int) {
			mouseDownChan <- button
		}, MouseUp: func(button int) {
			mouseUpChan <- button
		}, Wheel: func(amount int, hor bool) {
			mouseWheelChan <- amount
		},
		OnCursorMove: func() {

		}, OnCursorStop: func() {

		}, OnActivated: func() {
			onActivatedChan <- struct{}{}
		}, OnDeactivated: func() {
			onDeactivatedChan <- struct{}{}
		}, NormalKeyChan: normalKeyChan,
	}

	logic.SetConfig(
		def.Config{
			FunctionMap: def.FunctionMap{
				def.Activate: {
					IsAlt:   true,
					KeyCode: 74,
				},
				def.Deactivate:      {KeyCode: 71},
				def.MoveRight:       {KeyCode: 76},
				def.MoveUp:          {KeyCode: 75},
				def.MoveLeft:        {KeyCode: 72},
				def.MoveDown:        {KeyCode: 74},
				def.ClickLeft:       {KeyCode: 65},
				def.ClickRight:      {KeyCode: 68},
				def.ClickMiddle:     {KeyCode: 83},
				def.WheelUp:         {KeyCode: 85},
				def.WheelDown:       {KeyCode: 78},
				def.SniperMode:      {KeyCode: 32},
				def.TeleportForward: {KeyCode: 70},
			},
			DataMap: def.DataMap{
				def.CursorAccelerationH: "4.0",
				def.CursorAccelerationV: "4.0",
				def.CursorFrictionH:     "3.6",
				def.CursorFrictionV:     "3.6",
				def.WheelAccelerationH:  "40",
				def.WheelAccelerationV:  "40",
				def.WheelFrictionH:      "30",
				def.WheelFrictionV:      "30",
				def.SniperModeSpeedH:    "1",
				def.SniperModeSpeedV:    "1",
				def.TeleportDistanceH:   "300",
				def.TeleportDistanceV:   "300",
			},
		},
	)
	go logic.Loop()
}

func TestNormalKeyChan(t *testing.T) {
	a := assert.New(t)
	go func() {
		normalKeyChan <- 100
		a.Equal(uint32(0x55), <-normalKeyChan)
	}()
	time.Sleep(time.Millisecond)
	a.True(logic.OnKey(0x55, true))
	a.True(logic.OnKey(0x55, false))
}

func TestDeactivated(t *testing.T) {
	a := assert.New(t)
	a.False(logic.OnKey(71, true))
	a.False(logic.OnKey(71, false))
	a.False(logic.OnKey(74, true))
	a.False(logic.OnKey(74, false))
	a.False(logic.OnKey(76, true))
	a.False(logic.OnKey(76, false))
}

func TestActivateDeactivate(t *testing.T) {
	a := assert.New(t)

	a.False(logic.OnKey(alt, true))
	a.True(logic.OnKey(74, false))
	a.True(logic.OnKey(alt, false))
	a.True(logic.OnKey(74, false))

	a.True(logic.OnKey(71, true))
	a.True(logic.OnKey(71, false))
}
