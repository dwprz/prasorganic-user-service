package middleware

import (
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	conf   *config.Config
	logger *logrus.Logger
}

func New(conf *config.Config, l *logrus.Logger) *Middleware {
	return &Middleware{
		conf:   conf,
		logger: l,
	}
}
