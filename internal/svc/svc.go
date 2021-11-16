package svc

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

func Register() (err error) {
	b, err := exec.Command("where", "mouseable").Output()
	if err != nil {
		err = errors.Wrap(err, "where mouseable")
		return
	}

	err = exec.Command(
		"sc", "create", ServiceName,
		"binpath=", fmt.Sprintf("\"%s -run\"", b),
		"start=", "auto",
		"displayName=", "Mouseable Service",
	).Run()
	if ee, ok := err.(*exec.ExitError); ok {
		switch ee.ExitCode() {
		case 5:
			fmt.Println("Permission denied")
			return nil
		case 1073:
			fmt.Println("Already registered")
			return nil
		}
	}

	if err != nil {
		err = errors.Wrap(err, "sc create")
		return
	}

	return
}

func Unregister() (err error) {
	err = exec.Command("sc", "delete", ServiceName).Run()
	if ee, ok := err.(*exec.ExitError); ok {
		switch ee.ExitCode() {
		case 5:
			fmt.Println("Permission denied")
			return nil
		case 1060:
			fmt.Println("Not registered")
			return nil
		}
	}

	if err != nil {
		err = errors.Wrap(err, "sc delete")
		return
	}

	return
}

func Reload() (err error) {
	err = exec.Command("net", "stop", ServiceName).Run()
	if ee, ok := err.(*exec.ExitError); ok {
		switch ee.ExitCode() {
		case 2:
			fmt.Println("Service is not started, now it starts.")
			err = nil
		}
	}

	if err != nil {
		err = errors.Wrap(err, "net stop")
		return
	}

	err = exec.Command("net", "start", ServiceName).Run()
	if ee, ok := err.(*exec.ExitError); ok {
		switch ee.ExitCode() {

		}
	}

	if err != nil {
		err = errors.Wrap(err, "net start")
		return
	}

	return
}
