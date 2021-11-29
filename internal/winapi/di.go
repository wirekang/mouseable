package winapi

var DI struct {
	OnKey func(keyCode uint32, isDown bool) (preventDefault bool)
}
