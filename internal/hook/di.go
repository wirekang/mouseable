package hook

var DI struct {
	OnHook   func()
	OnUnhook func()
	OnKey    func(keyCode uint32, isDown bool) (preventDefault bool)
}
