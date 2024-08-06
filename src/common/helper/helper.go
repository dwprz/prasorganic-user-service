package helper

import (
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/interface/helper"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/sirupsen/logrus"
)

type HelperImpl struct {
	imageKit *imagekit.ImageKit
	conf     *config.Config
	logger   *logrus.Logger
}

func New(ik *imagekit.ImageKit, conf *config.Config, l *logrus.Logger) helper.Helper {
	return &HelperImpl{
		imageKit: ik,
		conf:     conf,
		logger:   l,
	}
}
