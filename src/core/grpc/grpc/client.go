package grpc

import (
	"github.com/dwprz/prasorganic-user-service/src/interface/client"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// this main grpc client
type Client struct {
	Otp      client.OtpGrpc
	otpConn *grpc.ClientConn
	logger   *logrus.Logger
}

func NewClient(ogc client.OtpGrpc, otpConn *grpc.ClientConn, l *logrus.Logger) *Client {

	return &Client{
		Otp:      ogc,
		otpConn: otpConn,
		logger:   l,
	}
}

func (g *Client) Close() {
	if err := g.otpConn.Close(); err != nil {
		g.logger.WithFields(logrus.Fields{"location": "grpc.Client/Close", "section": "otpConn.Close"}).Errorf(err.Error())
	}
}
