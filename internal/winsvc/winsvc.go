package winsvc

import (
	"context"
	"os"
	"strconv"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"

	"github.com/wirekang/mouseable/internal/lg"
	"github.com/wirekang/mouseable/internal/script"
)

type Handler struct {
	RunFunc func(ctx context.Context)
}

func (h Handler) Execute(
	args []string,
	requestChan <-chan svc.ChangeRequest,
	statusChan chan<- svc.Status,
) (svcSpecificEC bool, exitCode uint32) {
	elog, _ := eventlog.Open(script.ServiceName)
	f, err := os.Create(lg.FilePath)
	elog.Info(1, lg.FilePath)
	if err != nil {
		elog.Info(2, err.Error())
	} else {
		n, err := f.WriteString("asdf")
		if err != nil {
			elog.Info(3, err.Error())
		} else {
			elog.Info(4, strconv.Itoa(n))

		}
	}

	statusChan <- svc.Status{State: svc.StartPending}
	ctx, cancel := context.WithCancel(context.Background())
	go h.RunFunc(ctx)
	statusChan <- svc.Status{
		State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

LOOP:
	for {
		switch r := <-requestChan; r.Cmd {
		case svc.Stop, svc.Shutdown:
			lg.Logf("Receive Stop")
			cancel()
			break LOOP
		case svc.Interrogate:
			statusChan <- r.CurrentStatus
			time.Sleep(100 & time.Millisecond)
			statusChan <- r.CurrentStatus
		}
	}

	statusChan <- svc.Status{State: svc.StopPending}
	return
}
