package logic

var DI struct {
	SetCursorPos func(x, y int)
	AddCursorPos func(dx, dy int32)
	GetCursorPos func() (x, y int)
	MouseDown    func(button int)
	MouseUp      func(button int)
	Wheel        func(amount int, hor bool)
}
