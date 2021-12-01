package logic

var DI struct {
	SetCursorPos  func(x, y int)
	AddCursorPos  func(dx, dy int)
	GetCursorPos  func() (x, y int)
	MouseDown     func(button int)
	MouseUp       func(button int)
	Wheel         func(amount int, hor bool)
	OnCursorMove  func()
	OnCursorStop  func()
	OnActivated   func()
	OnDeactivated func()
	NormalKeyChan chan uint32
}
