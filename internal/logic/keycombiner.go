package logic

import (
	"strings"
	"time"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/typ"
)

type keyCombiner struct {
	preventKeyUpMap  map[typ.Key]struct{}
	pressingModKey   typ.Key
	lastDownKey      typ.Key
	lastDownKeyTime  int64
	doublePressSpeed int64

	keyInfoChan        <-chan typ.KeyInfo
	preventDefaultChan chan<- bool

	needNextKeyChan <-chan struct{}
	nextKeyChan     chan<- typ.Key

	internalKeyInfoChan        chan<- typ.KeyInfo
	internalPreventDefaultChan <-chan bool

	configChan <-chan typ.Config
}

func (k *keyCombiner) Run() {
	lg.Printf("Run keyCombiner")
	for {
		select {
		case config := <-k.configChan:
			k.changeConfig(config)
		case ki := <-k.keyInfoChan:
			k.proc(ki.Key, ki.IsDown)
		}
	}
}

func (k *keyCombiner) changeConfig(config typ.Config) {
	k.doublePressSpeed = int64(config.DataValue("double-press-speed").Int())
}

func (k *keyCombiner) proc(originKey typ.Key, isDown bool) {
	combinedKey, isMod := normalizeModKey(originKey)

	if isMod {
		k.setPressingModKey(combinedKey, isDown)
	}

	preventDefault := false
	if isDown {
		if !isMod {
			combinedKey = k.combineModKey(combinedKey)
		}

		if k.isDouble(originKey) {
			combinedKey += "x2"
		}

		k.internalKeyInfoChan <- typ.KeyInfo{
			Key:    combinedKey,
			IsDown: isDown,
		}

		preventDefault = <-k.internalPreventDefaultChan || preventDefault
		preventDefault = k.isNeedNextKey(combinedKey) || preventDefault
	} else {
		preventDefault = k.popPreventKeyUpMap(originKey) || preventDefault
	}

	k.preventDefaultChan <- preventDefault
}

func (k *keyCombiner) setPressingModKey(key typ.Key, isDown bool) {
	if isDown {
		k.pressingModKey = key
	} else if k.pressingModKey == key {
		k.pressingModKey = ""
	}
}

// combineModKey convert "A" to "Shift+A" if Shift was pressing.
func (k *keyCombiner) combineModKey(key typ.Key) (rst typ.Key) {
	pressingModKey := k.pressingModKey
	if pressingModKey != "" {
		rst = pressingModKey + "+"
	}
	rst += key
	return
}

// isDouble returns true if key was pressed recently.
func (k *keyCombiner) isDouble(key typ.Key) (ok bool) {
	lastTime := k.lastDownKeyTime
	if k.lastDownKey != key {
		return
	}
	ok = time.Now().UnixMilli()-lastTime <= k.doublePressSpeed
	k.lastDownKeyTime = lastTime
	k.lastDownKey = key
	return
}

func (k *keyCombiner) isNeedNextKey(editedKey typ.Key) bool {
	select {
	case <-k.needNextKeyChan:
		k.nextKeyChan <- editedKey
		return true
	default:
		return false
	}
}

func (k *keyCombiner) popPreventKeyUpMap(key typ.Key) (ok bool) {
	_, ok = k.preventKeyUpMap[key]
	if ok {
		delete(k.preventKeyUpMap, key)
	}
	return
}

// normalizeModKey converts "Left Shift" to "Shift", "Right Ctrl" to "Ctrl". Otherwise, return as it is
func normalizeModKey(key typ.Key) (v typ.Key, isModKey bool) {
	v = key
	isModKey = strings.Contains(string(key), "Alt")
	if isModKey {
		v = "Alt"
		return
	}

	isModKey = strings.Contains(string(key), "Shift")
	if isModKey {
		v = "Shift"
		return
	}

	isModKey = strings.Contains(string(key), "Ctrl")
	if isModKey {
		v = "Ctrl"
		return
	}
	return
}
