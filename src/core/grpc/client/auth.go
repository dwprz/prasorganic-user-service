package client

import (
	"context"
	"fmt"
	"log"

	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/interface/client"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OtpGrpcImpl struct {
	client   pb.OtpServiceClient
	cbreaker *gobreaker.CircuitBreaker[any]
}

func NewOtpGrpc(cb *gobreaker.CircuitBreaker[any], conf *config.Config, unaryRequest *interceptor.UnaryRequest) (client.OtpGrpc, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryRequest.AddBasicAuth),
	)

	conn, err := grpc.NewClient(conf.ApiGateway.BaseUrl, opts...)
	if err != nil {
		log.Fatalf("new otp grpc client: %v", err.Error())
	}

	client := pb.NewOtpServiceClient(conn)

	return &OtpGrpcImpl{
		client:   client,
		cbreaker: cb,
	}, conn
}

func (u *OtpGrpcImpl) Send(ctx context.Context, email string) error {
	_, err := u.cbreaker.Execute(func() (any, error) {
		_, err := u.client.Send(ctx, &pb.SendRequest{
			Email: email,
		})
		return nil, err
	})

	return err
}

func (u *OtpGrpcImpl) Verify(ctx context.Context, data *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	res, err := u.cbreaker.Execute(func() (any, error) {
		res, err := u.client.Verify(ctx, &pb.VerifyRequest{
			Email: data.Email,
			Otp:   data.Otp,
		})
		return res, err
	})

	if err != nil {
		return nil, err
	}

	user, ok := res.(*pb.VerifyResponse)
	if !ok {
		return nil, fmt.Errorf("client.OtpGrpcImpl/Verify | unexpected type: %T", res)
	}

	return user, err
}
