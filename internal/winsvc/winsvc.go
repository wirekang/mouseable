package winsvc

import (
	"context"
	"time"

	"golang.org/x/sys/windows/svc"

	"github.com/wirekang/mouseable/internal/lg"
)

type Handler struct {
	RunFunc func(ctx context.Context)
}

func (h Handler) Execute(
	args []string,
	requestChan <-chan svc.ChangeRequest,
	statusChan chan<- svc.Status,
) (svcSpecificEC bool, exitCode uint32) {
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
