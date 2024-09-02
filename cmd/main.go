package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc"
	"github.com/dwprz/prasorganic-user-service/src/core/restful"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/src/service"
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

	postgresDB := database.NewPostgres()
	defer database.ClosePostgres(postgresDB)

	redisDB := database.NewRedisCluster()
	defer redisDB.Close()

	userCache := cache.NewUser(redisDB)
	userRepository := repository.NewUser(postgresDB, userCache)

	grpcClient := grpc.InitClient()
	defer grpcClient.Close()

	userService := service.NewUser(grpcClient, userRepository, userCache)

	grpcServer := grpc.InitServer(userService)
	defer grpcServer.Stop()

	go grpcServer.Run()

	restfulClient := restful.InitClient()
	restfulServer := restful.InitServer(restfulClient, userService)
	defer restfulServer.Stop()

	go restfulServer.Run()

	<-closeCH
}
