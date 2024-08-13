package restful

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HandleJsonError(c *fiber.Ctx, err *json.UnmarshalTypeError) error {
	return c.Status(400).JSON(fiber.Map{
		"field":  err.Field,
		"errros": fmt.Sprintf("%s must be a %v not an %s", err.Field, err.Type, err.Value),
	})
}
