package task

import (
	"goboot/pkg/errorx"
	"goboot/pkg/log"
	"goboot/pkg/queue"
	"time"
)

func showDataLenInQueue(name string, logger *log.Logger, queue queue.Queue, max int) {
	go func() {
		defer errorx.Recover(logger)
		for {
			logger.Warnf("task count：%-6d, name: %-30s", queue.Len(), name)
			time.Sleep(time.Second * 10)
		}
	}()
}
