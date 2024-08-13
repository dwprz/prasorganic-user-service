package server

import (
	"fmt"
	"net"

	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Grpc struct {
	port                     string
	server                   *grpc.Server
	userGrpcHandler          user.UserServiceServer
	unaryResponseInterceptor *interceptor.UnaryResponse
}

// this main grpc server
func NewGrpc(userGrpcHandler user.UserServiceServer, uri *interceptor.UnaryResponse) *Grpc {
	port := config.Conf.CurrentApp.GrpcPort

	return &Grpc{
		port:                     port,
		userGrpcHandler:          userGrpcHandler,
		unaryResponseInterceptor: uri,
	}
}

func (g *Grpc) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", g.port))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "server.Grpc/Run", "section": "net.Listen"}).Fatal(err)
	}

	log.Logger.Infof("grpc run in port: %s", g.port)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			g.unaryResponseInterceptor.Recovery,
			g.unaryResponseInterceptor.Error,
		))

	g.server = grpcServer

	user.RegisterUserServiceServer(grpcServer, g.userGrpcHandler)

	if err := grpcServer.Serve(listener); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "server.Grpc/Run", "section": "grpcServer.Serve"}).Fatal(err)
	}
}

func (g *Grpc) Stop() {
	g.server.Stop()
}
