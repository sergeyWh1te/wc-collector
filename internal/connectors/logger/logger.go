package logger

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger            *logrus.Logger
	onceDefaultClient sync.Once
)

func New(f *os.File) *logrus.Logger {
	onceDefaultClient.Do(func() {
		logger = logrus.StandardLogger()

		wrt := io.MultiWriter(os.Stdout, f)
		logger.SetOutput(wrt)
	})

	return logger
}
