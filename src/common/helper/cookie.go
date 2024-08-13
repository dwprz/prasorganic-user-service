package helper

import (
	"time"
	"github.com/gofiber/fiber/v2"
)

func ClearCookie(name string, path string) *fiber.Cookie {
	clearCookie := &fiber.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		HTTPOnly: true,
		Expires:  time.Now().Add(-time.Hour),
	}

	return clearCookie
}
