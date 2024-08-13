package restful

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func HandleStrconvError(c *fiber.Ctx, err *strconv.NumError) error {
	return c.Status(400).JSON(fiber.Map{
		"errors": fmt.Sprintf("%v (%s)", err.Err, err.Num),
	})
}
