package middleware

import (
	"context"

	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (m *Middleware) Error(c *fiber.Ctx, err error) error {
	m.logger.WithFields(logrus.Fields{
		"host":     c.Hostname(),
		"ip":       c.IP(),
		"protocol": c.Protocol(),
		"location": c.OriginalURL(),
		"method":   c.Method(),
		"from":     "error middleware",
	}).Error(err.Error())

	if err != nil && c.OriginalURL() == "/api/users/current/photo-profile" && c.Method() == "PATCH" {
		filename := c.Locals("filename").(string)
		if filename != "" {
			go m.helper.DeleteFile("./tmp/" + filename)
		}

		req, ok := c.Locals("update_photo_profile_req").(dto.UpdatePhotoProfileReq)
		if ok && req.NewPhotoProfileId != "" {
			go m.imageKit.Media.DeleteFile(context.Background(), req.NewPhotoProfileId)
		}
	}

	if validationError, ok := err.(validator.ValidationErrors); ok {

		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]any{
				"field":       validationError[0].Field(),
				"description": validationError[0].Error(),
			},
		})
	}

	if responseError, ok := err.(*errors.Response); ok {
		return c.Status(int(responseError.HttpCode)).JSON(fiber.Map{
			"errors": responseError.Message,
		})
	}

	return c.Status(500).JSON(fiber.Map{
		"errors": "sorry, internal server error try again later",
	})
}
