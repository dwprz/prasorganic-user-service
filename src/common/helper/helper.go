package helper

import (
	"github.com/dwprz/prasorganic-user-service/src/interface/helper"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/sirupsen/logrus"
)

type HelperImpl struct {
	conf   *config.Config
	logger *logrus.Logger
}

func New(conf *config.Config, l *logrus.Logger) helper.Helper {
	return &HelperImpl{
		conf:   conf,
		logger: l,
	}
}
