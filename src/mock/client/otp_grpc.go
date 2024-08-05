package client

import (
	"context"

	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/stretchr/testify/mock"
)

type OtpGrpc struct {
	mock.Mock
}

func NewOtpGrpcMock() *OtpGrpc {
	return &OtpGrpc{
		Mock: mock.Mock{},
	}
}

func (o *OtpGrpc) Send(ctx context.Context, email string) error {
	arguments := o.Mock.Called(ctx, email)

	return arguments.Error(0)
}

func (o *OtpGrpc) Verify(ctx context.Context, data *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	arguments := o.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*pb.VerifyResponse), arguments.Error(1)
}
