package util

import (
	"log"

	userpb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	grpcapp "github.com/dwprz/prasorganic-user-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func NewGrpcServer() (*grpcapp.Server, *gorm.DB, *redis.ClusterClient, *config.Config, *logrus.Logger) {
	logger := logger.New()
	validator := validator.New()
	conf := config.New("DEVELOPMENT", logger)

	postgresDB := database.NewPostgres(conf)
	redisDB := database.NewRedisCluster(conf)

	userCache := cache.NewUser(redisDB, logger)
	userRepository := repository.NewUser(postgresDB, userCache)
	userService := service.NewUser(validator, userRepository, userCache)
	unaryRespInterceptor := interceptor.NewUnaryResponse(logger)

	userGrpcServer := server.NewUserGrpc(logger, userService)
	grpcServer := grpcapp.NewServer(conf.CurrentApp.GrpcPort, userGrpcServer, unaryRespInterceptor, logger)

	return grpcServer, postgresDB, redisDB, conf, logger
}

func NewGrpcUserClient(apiGatewayBaseUrl string) (userpb.UserServiceClient, *grpc.ClientConn) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(apiGatewayBaseUrl, opts...)

	if err != nil {
		log.Fatal("failed to create new grpc user client")
	}

	userServiceClient := userpb.NewUserServiceClient(conn)

	return userServiceClient, conn
}