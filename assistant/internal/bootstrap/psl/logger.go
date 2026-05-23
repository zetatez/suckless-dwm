package psl

import (
	"fmt"
	"sync"

	"assistant/pkg/xlog"

	"github.com/sirupsen/logrus"
)

var (
	logger  *logrus.Logger
	onceLog sync.Once
)

func GetLogger() *logrus.Logger {
	return logger
}

func InitLog() error {
	var initErr error
	onceLog.Do(func() {
		var err error
		logger, err = xlog.NewLogger(GetConfig().Log)
		if err != nil {
			initErr = fmt.Errorf("new logger failed: %w", err)
			return
		}
	})
	return initErr
}
