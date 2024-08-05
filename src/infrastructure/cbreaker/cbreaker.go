package cbreaker

import (
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
)

type CircuitBreaker struct {
	OtpGrpc *gobreaker.CircuitBreaker[any]
}

func New(logger *logrus.Logger) *CircuitBreaker {
	otpGrpcCBreaker := setupForOtpGrpc(logger)

	return &CircuitBreaker{
		OtpGrpc: otpGrpcCBreaker,
	}
}