package imagekit

import (
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/logger"
)

func New(conf *config.Config) *imagekit.ImageKit {
	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  conf.ImageKit.PrivateKey,
		PublicKey:   conf.ImageKit.PublicKey,
		UrlEndpoint: conf.ImageKit.BaseUrl,
	})

	ik.Logger.SetLevel(logger.ERROR)

	return ik
}
