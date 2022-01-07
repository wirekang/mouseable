package io

import (
	"os"

	"github.com/juju/fslock"
)

var lock *fslock.Lock

func Lock() (ok bool) {
	lockFile := configDir + "\\lockfile"
	_ = os.Mkdir(configDir, os.ModeDir)
	lock = fslock.New(lockFile)
	err := lock.TryLock()
	if err != nil {
		return
	}
	ok = true
	return
}

func Unlock() {
	lock.Unlock()
}
