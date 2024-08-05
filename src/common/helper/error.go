package helper

import (
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

		c.Status(500).JSON(fiber.Map{
			"errors": "sorry, internal server error try again later",
		})
	}
}
