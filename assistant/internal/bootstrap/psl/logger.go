package psl

import (
	"sync"

	"assistant/pkg/xlog"

	"github.com/sirupsen/logrus"
)

var (
	logger     *logrus.Logger
	onceLogger sync.Once
)

func GetLogger() *logrus.Logger { return logger }

func InitLog() error {
	var err error
	onceLogger.Do(func() {
		logger, err = xlog.NewLogger(GetConfig().Log)
	})
	return err
}
