package logic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/wirekang/mouseable/internal/config"
)

func Loop(ctx context.Context) (err error) {
	cfg, err := config.Load()
	_ = cfg
	if err != nil {
		err = errors.Wrap(err, "config.Load")
		return
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP

			// todo
			// 키보드 입력을 받아서 상태를 변화시킨다.
		}

		// 상태에 따라서 마우스 커서를 이동시킨다.
	}
	return
}

// winsvc 패키지에서 키보드 입력을 처리하고 Output은 key 패키지에서 처리된 string으로 반환한다.
