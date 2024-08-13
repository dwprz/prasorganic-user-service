package client

import "github.com/dwprz/prasorganic-user-service/src/interface/delivery"

// this main restful client
type Restful struct {
	ImageKit delivery.ImageKitRestful
}

func NewRestful(ikd delivery.ImageKitRestful) *Restful {
	return &Restful{
		ImageKit: ikd,
	}
}
