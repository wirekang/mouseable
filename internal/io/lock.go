package io

import (
	"github.com/juju/fslock"

	"github.com/wirekang/mouseable/internal/must"
)

var lock *fslock.Lock

func Lock() (ok bool) {
	lockFile := must.ConfigDir() + "\\lockfile"
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
