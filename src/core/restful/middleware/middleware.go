package middleware

import (
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/interface/helper"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	imageKit *imagekit.ImageKit
	conf     *config.Config
	helper   helper.Helper
	logger   *logrus.Logger
}

func New(ik *imagekit.ImageKit, conf *config.Config, h helper.Helper, l *logrus.Logger) *Middleware {
	return &Middleware{
		imageKit: ik,
		conf:     conf,
		helper:   h,
		logger:   l,
	}
}
