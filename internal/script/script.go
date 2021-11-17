package script

import (
	"os/exec"

	"github.com/wirekang/winsvc/internal/config"
)

func OpenConfigFile() (err error) {
	_ = exec.Command(
		"explorer", config.FilePath,
	).Run()
	return
}
