package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func SetupLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return log
}
