package restful

import (
	"net/http"
	"time"

	"github.com/dwprz/prasorganic-user-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/router"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// this main restful server
type Server struct {
	app                *fiber.App
	userRestfulHandler *handler.UserRestful
	middleware         *middleware.Middleware
	conf               *config.Config
}

func NewServer(urh *handler.UserRestful, m *middleware.Middleware, conf *config.Config) *Server {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		IdleTimeout:   20 * time.Second,
		ReadTimeout:   20 * time.Second,
		WriteTimeout:  20 * time.Second,
		ErrorHandler:  m.Error,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://restful.local:80",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	router.Create(app, urh, m)

	return &Server{
		app:                app,
		userRestfulHandler: urh,
		middleware:         m,
		conf:               conf,
	}
}

func (r *Server) Run() {
	r.app.Listen(r.conf.CurrentApp.RestfulAddress)
}

func (r *Server) Test(req *http.Request) (*http.Response, error) {
	res, err := r.app.Test(req)

	return res, err
}

func (r *Server) Stop() {
	r.app.Shutdown()
}
