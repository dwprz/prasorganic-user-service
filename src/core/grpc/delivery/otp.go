package delivery

import (
	"context"
	"fmt"

	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/cbreaker"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/interface/delivery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OtpGrpcImpl struct {
	client   pb.OtpServiceClient
}

func NewOtpGrpc(unaryRequest *interceptor.UnaryRequest) (delivery.OtpGrpc, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryRequest.AddBasicAuth),
	)

	conn, err := grpc.NewClient(config.Conf.ApiGateway.BaseUrl, opts...)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.NewOtpGrpc", "section": "grpc.NewClient"}).Fatal(err)
	}

	client := pb.NewOtpServiceClient(conn)

	return &OtpGrpcImpl{
		client:   client,
	}, conn
}

func (u *OtpGrpcImpl) Send(ctx context.Context, email string) error {
	_, err := cbreaker.OtpGrpc.Execute(func() (any, error) {
		_, err := u.client.Send(ctx, &pb.SendReq{
			Email: email,
		})
		return nil, err
	})

	return err
}

func (u *OtpGrpcImpl) Verify(ctx context.Context, data *pb.VerifyReq) (*pb.VerifyRes, error) {
	res, err := cbreaker.OtpGrpc.Execute(func() (any, error) {
		res, err := u.client.Verify(ctx, &pb.VerifyReq{
			Email: data.Email,
			Otp:   data.Otp,
		})
		return res, err
	})

	if err != nil {
		return nil, err
	}

	user, ok := res.(*pb.VerifyRes)
	if !ok {
		return nil, fmt.Errorf("client.OtpGrpcImpl/Verify | unexpected type: %T", res)
	}

	return user, err
}
