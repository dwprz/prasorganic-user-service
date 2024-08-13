package middleware

import (
	"github.com/dwprz/prasorganic-user-service/src/common/helper"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) SaveTemporaryImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	var maxSize int64 = 1 * 1000 * 1000 // 1 mb
	if file.Size > maxSize {
		return c.Status(400).JSON(fiber.Map{"errors": "file size is too large"})
	}

	filename := helper.CreateUnixFileName(file.Filename)
	path := "./tmp/" + filename

	if err := helper.CheckExistDir("./tmp"); err != nil {
		return err
	}

	c.SaveFile(file, path)

	c.Locals("filename", filename)
	return c.Next()
}
