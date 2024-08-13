package imagekit

import (
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/logger"
)

var IK *imagekit.ImageKit

func init() {
	IK = imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  config.Conf.ImageKit.PrivateKey,
		PublicKey:   config.Conf.ImageKit.PublicKey,
		UrlEndpoint: config.Conf.ImageKit.BaseUrl,
	})

	IK.Logger.SetLevel(logger.ERROR)
}
