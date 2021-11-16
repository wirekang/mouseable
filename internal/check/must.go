package check

import (
	"runtime"
)

func MustWindows() {
	if runtime.GOOS != "windows" {
		panic("not windows")
	}
}
