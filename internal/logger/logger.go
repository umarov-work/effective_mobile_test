package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	return log
}
