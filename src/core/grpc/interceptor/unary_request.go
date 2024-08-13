package interceptor

import (
	"context"
	"encoding/base64"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UnaryRequest struct {}

func NewUnaryRequest() *UnaryRequest {
	return &UnaryRequest{}
}

func (u *UnaryRequest) AddBasicAuth(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(config.Conf.ApiGateway.BasicAuth))
	md.Append("Authorization", "Basic "+auth)

	authCtx := metadata.NewOutgoingContext(ctx, md)

	return invoker(authCtx, method, req, reply, cc, opts...)
}