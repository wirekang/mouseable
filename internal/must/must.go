package must

import (
	"os"
	"runtime"
	"strings"
)

func Windows() {
	if runtime.GOOS != "windows" {
		panic("not windows")
	}
}

func ConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(homeDir, "/", "\\")
}
