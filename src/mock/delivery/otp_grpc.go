package delivery

import (
	"context"

	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/stretchr/testify/mock"
)

type OtpGrpcMock struct {
	mock.Mock
}

func NewOtpGrpcMock() *OtpGrpcMock {
	return &OtpGrpcMock{
		Mock: mock.Mock{},
	}
}

func (o *OtpGrpcMock) Send(ctx context.Context, email string) error {
	arguments := o.Mock.Called(ctx, email)

	return arguments.Error(0)
}

func (o *OtpGrpcMock) Verify(ctx context.Context, data *pb.VerifyReq) (*pb.VerifyRes, error) {
	arguments := o.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*pb.VerifyRes), arguments.Error(1)
}
