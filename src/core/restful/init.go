package restful

import (
	"github.com/dwprz/prasorganic-user-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/delivery"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
)

func InitServer(rc *client.Restful, us service.User) *server.Restful {

	userHandler := handler.NewUser(us, rc)
	middleware := middleware.New(rc)

	restfulServer := server.NewRestful(userHandler, middleware)
	return restfulServer
}

func InitClient() *client.Restful {
	imageKitDelivery := delivery.NewImageKit()
	
	restfulClient := client.NewRestful(imageKitDelivery)
	return restfulClient
}
