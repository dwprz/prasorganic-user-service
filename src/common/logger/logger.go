package logger

import "github.com/sirupsen/logrus"

func New() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&customJSONFormatter{})

	return logger
}
