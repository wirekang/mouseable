package hook

import (
	"os"
	"os/signal"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"github.com/pkg/errors"

	"github.com/wirekang/winsvc/internal/lg"
)

func Loop() (err error) {
	keyboardChan := make(chan types.KeyboardEvent, 100)

	err = keyboard.Install(nil, keyboardChan)
	if err != nil {
		err = errors.Wrap(err, "keyboard.Install")
		return
	}

	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	lg.Logf("Start capturing keyboard input")

	for {
		select {
		case <-signalChan:
			lg.Logf("Received shutdown signal")
			return nil
		case k := <-keyboardChan:
			lg.Logf("Message: '%v'     VKCode: '%v'\n", k.Message, k.VKCode)
		}
	}
}
