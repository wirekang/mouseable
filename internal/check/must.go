package check

import (
	"os"
	"runtime"
	"strings"
)

func MustWindows() {
	if runtime.GOOS != "windows" {
		panic("not windows")
	}
}

func MustCacheDir() string {
	d, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	dir := strings.ReplaceAll(d, "\\", "/") + "/mouseable"
	_ = os.Mkdir(dir, 0777)
	return dir
}
