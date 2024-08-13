package client

import (
	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/dwprz/prasorganic-user-service/src/interface/delivery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// this main grpc client
type Grpc struct {
	Otp      delivery.OtpGrpc
	otpConn *grpc.ClientConn
}

func NewGrpc(ogd delivery.OtpGrpc, otpConn *grpc.ClientConn) *Grpc {

	return &Grpc{
		Otp:      ogd,
		otpConn: otpConn,
	}
}

func (g *Grpc) Close() {
	if err := g.otpConn.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "client.Grpc/Close", "section": "otpConn.Close"}).Error(err)
	}
}
