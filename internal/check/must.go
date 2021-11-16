package check

import (
	"os"
	"runtime"
)

func MustWindows() {
	if runtime.GOOS != "windows" {
		panic("not windows")
	}
}

func MustConfigDir() string {
	d, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir := d[:2] + "\\mouseable\\"
	_ = os.Mkdir(dir, 0777)
	return dir
}
