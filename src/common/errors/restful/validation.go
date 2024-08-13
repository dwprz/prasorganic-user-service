package restful

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func HandleValidationError(c *fiber.Ctx, err validator.ValidationErrors) error {
	return c.Status(400).JSON(fiber.Map{
		"errors": map[string]any{
			"field":       err[0].Field(),
			"description": err[0].Error(),
		},
	})
}
