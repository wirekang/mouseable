package script

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/check"
)

func Register() (err error) {
	b, err := exec.Command("where", "mouseable").Output()
	if err != nil {
		err = errors.Wrap(err, "where mouseable")
		return
	}

	err = exec.Command(
		"sc", "create", ServiceName,
		"binpath=", strings.TrimSpace(string(b))+" -run",
		"start=", "auto",
		"displayName=", "Mouseable Service",
	).Run()

	if checkExitCode(err) {
		return
	}

	if err != nil {
		err = errors.Wrap(err, "sc create")
		return
	}

	err = exec.Command("sc", "start", ServiceName).Run()

	if checkExitCode(err) {
		return
	}

	if err != nil {
		err = errors.Wrap(err, "sc start")
		return
	}

	return
}

func Unregister() (err error) {
	_ = exec.Command("net", "stop", ServiceName).Run()
	_ = exec.Command("taskkill", "/F", "/IM", "mmc.exe").Run()
	err = exec.Command("sc", "delete", ServiceName).Run()

	if checkExitCode(err) {
		return
	}

	if err != nil {
		err = errors.Wrap(err, "sc delete")
		return
	}

	return
}

func Reload() (err error) {
	_ = exec.Command("net", "stop", ServiceName).Run()

	err = exec.Command("sc", "start", ServiceName).Run()

	if checkExitCode(err) {
		return
	}

	if err != nil {
		err = errors.Wrap(err, "sc start")
		return
	}

	return
}

func OpenConfigDir() (err error) {
	_ = exec.Command(
		"explorer", check.MustConfigDir(),
	).Run()
	return
}

func checkExitCode(err error) bool {
	if ee, ok := (err).(*exec.ExitError); ok {
		fmt.Println("\n------[maybe help]------")
		fmt.Println()
		switch ee.ExitCode() {
		case 5:
			fmt.Println("Permission denied")
		case 1060:
			fmt.Println("Not registered")
		case 1073:
			fmt.Println("Already registered")
		default:
			fmt.Println("Nothing to show")
		}
		fmt.Println("\n------------------------")
		fmt.Println()
		return true
	}
	return false
}
