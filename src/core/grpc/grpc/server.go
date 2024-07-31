package grpc

import (
	"fmt"
	"net"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	port                     string
	server                   *grpc.Server
	userServiceServer        user.UserServiceServer
	unaryResponseInterceptor *interceptor.UnaryResponse
	logger                   *logrus.Logger
}

// this main grpc server
func NewServer(port string, uss user.UserServiceServer, uri *interceptor.UnaryResponse, l *logrus.Logger) *Server {
	return &Server{
		port:                     port,
		userServiceServer:        uss,
		unaryResponseInterceptor: uri,
		logger:                   l,
	}
}

func (g *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", g.port))
	if err != nil {
		g.logger.Errorf("failed to listen on port %s : %v", g.port, err)
	}

	g.logger.Infof("grpc run in port: %s", g.port)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			g.unaryResponseInterceptor.Recovery,
			g.unaryResponseInterceptor.Error,
		))

	g.server = grpcServer

	user.RegisterUserServiceServer(grpcServer, g.userServiceServer)

	if err := grpcServer.Serve(listener); err != nil {
		g.logger.Errorf("failed to serve grpc on port %s : %v", g.port, err)
	}
}

func (g *Server) Stop() {
	g.server.Stop()
}
