package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

const (
	AppName = "todolist"
)

func InitLog(logLevel string) error {

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	logrus.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(logLevel)

	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
		return err
	}

	logrus.SetLevel(level)
	return nil
}
