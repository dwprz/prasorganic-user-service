package restful

import (
	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LogError(c *fiber.Ctx, err error) {
	log.Logger.WithFields(logrus.Fields{
		"host":     c.Hostname(),
		"ip":       c.IP(),
		"protocol": c.Protocol(),
		"location": c.OriginalURL(),
		"method":   c.Method(),
		"from":     "error middleware",
	}).Error(err)
}
