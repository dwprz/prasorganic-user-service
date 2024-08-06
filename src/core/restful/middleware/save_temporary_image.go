package middleware

import (
	"fmt"
	"regexp"
	"time"

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

	re := regexp.MustCompile(`[ %?#&=]`)
	encodedName := re.ReplaceAllString(file.Filename, "-")
	epochMillis := time.Now().UnixMilli()

	filename := fmt.Sprintf("%d-%s", epochMillis, encodedName)

	c.SaveFile(file, "./tmp/"+filename)

	c.Locals("filename", filename)
	return c.Next()
}
