package helper

import (
	"context"

	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *HelperImpl) HandlePanic(c *fiber.Ctx) {
	message := recover()

	if message != nil {
		h.logger.WithFields(logrus.Fields{
			"host":     c.Hostname(),
			"ip":       c.IP(),
			"protocol": c.Protocol(),
			"location": c.OriginalURL(),
			"from":     "Handle Panic",
		}).Error(message)

		if c.OriginalURL() == "/api/users/current/photo-profile" && c.Method() == "PATCH" {
			filename := c.Locals("filename").(string)
			if filename != "" {
				go h.DeleteFile("./tmp/" + filename)
			}

			req, ok := c.Locals("update_photo_profile_req").(dto.UpdatePhotoProfileReq)
			if ok && req.NewPhotoProfileId != "" {
				go h.imageKit.Media.DeleteFile(context.Background(), req.NewPhotoProfileId)
			}
		}

		c.Status(500).JSON(fiber.Map{
			"errors": "sorry, internal server error try again later",
		})
	}
}
