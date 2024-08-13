package router

import (
	"github.com/dwprz/prasorganic-user-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/middleware"
	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App, h *handler.User, m *middleware.Middleware) {
	app.Add("GET", "/api/users/current", m.VerifyJwt, h.GetCurrent)
	app.Add("PATCH", "/api/users/current", m.VerifyJwt, h.UpdateProfile)
	app.Add("PATCH", "/api/users/current/password", m.VerifyJwt, h.UpdatePassword)
	app.Add("PATCH", "/api/users/current/email", m.VerifyJwt, h.UpdateEmail)
	app.Add("PATCH", "/api/users/current/email/verify", m.VerifyJwt, h.VerifyUpdateEmail)
	app.Add("PATCH", "/api/users/current/photo-profile", m.VerifyJwt, m.SaveTemporaryImage, m.ValidateImage, m.UploadToImageKit, h.UpdatePhotoProfile)
}
