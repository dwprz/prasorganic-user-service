package helper

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/gofiber/fiber/v2"
)

type Helper interface {
	GenerateAccessToken(userId string, email string, role string) (string, error)
	GetMetadata(ctx context.Context) *entity.Metadata 
	DeleteFile(path string)
	HandlePanic(c *fiber.Ctx)
	ClearCookie(name string, path string) *fiber.Cookie
}
