package router

import (
	"github.com/dwprz/prasorganic-user-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/middleware"
	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App, h *handler.UserRestful, m *middleware.Middleware) {
	app.Add("GET", "/api/users/current", m.VerifyJwt, h.GetCurrent)
}
