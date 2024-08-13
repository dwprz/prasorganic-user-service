package middleware

import (
	"github.com/dwprz/prasorganic-user-service/src/common/helper"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) UploadToImageKit(c *fiber.Ctx) error {
	filename := c.Locals("filename").(string)
	path := "./tmp/" + filename

	res, err := m.restfulClient.ImageKit.UploadImage(c.Context(), path, filename)
	if err != nil {
		return err
	}

	c.Locals("upload_imagekit_result", res)
	go helper.DeleteFile("./tmp/" + filename)

	return c.Next()
}
