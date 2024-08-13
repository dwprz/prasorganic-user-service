package grpc

import (
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/delivery"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/handler"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
)

func InitServer(us service.User) *server.Grpc {
	userHandler := handler.NewUserGrpc(us)
	unaryResponseInterceptor := interceptor.NewUnaryResponse()

	grpcServer := server.NewGrpc(userHandler, unaryResponseInterceptor)
	return grpcServer
}

func InitClient() *client.Grpc {
	unaryRequestInterceptor := interceptor.NewUnaryRequest()
	otpDelivery, otpConn := delivery.NewOtpGrpc(unaryRequestInterceptor)

	grpcClient := client.NewGrpc(otpDelivery, otpConn)
	return grpcClient
}
