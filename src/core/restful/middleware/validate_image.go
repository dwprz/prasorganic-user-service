package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/filetype"
	"os"
)

// validasi image berdasarkan magic number
func (m *Middleware) ValidateImage(c *fiber.Ctx) error {
	filename := c.Locals("filename").(string)

	fileCeck, err := os.Open("./tmp/" + filename)
	if err != nil {
		return err
	}
	defer fileCeck.Close()

	head := make([]byte, 261)
	fileCeck.Read(head)

	if !filetype.IsImage(head) {
		return c.Status(400).JSON(fiber.Map{"errors": "file is not an image"})
	}

	return c.Next()
}
