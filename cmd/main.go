package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/restful"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/src/service"
	"github.com/go-playground/validator/v10"
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func handleCloseApp(closeCH chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		close(closeCH)
	}()
}

func main() {
	closeCH := make(chan struct{})
	handleCloseApp(closeCH)

	appStatus := os.Getenv("PRASORGANIC_APP_STATUS")

	logger := logger.New()
	conf := config.New(appStatus, logger)
	postgresDb := database.NewPostgres(conf)
	redisDb := database.NewRedisCluster(conf)
	validator := validator.New()

	userCache := cache.NewUser(redisDb, logger)
	userRepository := repository.NewUser(postgresDb, userCache)
	userService := service.NewUser(validator, userRepository, userCache)
	unaryResInterceptor := interceptor.NewUnaryResponse(logger)
	userGrpcServer := server.NewUserGrpc(logger, userService)

	grpcServer := grpc.NewServer(conf.CurrentApp.GrpcPort, userGrpcServer, unaryResInterceptor, logger)
	defer grpcServer.Stop()

	go grpcServer.Run()

	userRestfulHandler := handler.NewUserRestful(userService)
	middleware := middleware.New(conf, logger)

	restfulServer := restful.NewServer(userRestfulHandler, middleware, conf)
	defer restfulServer.Stop()

	go restfulServer.Run()

	<-closeCH
}
