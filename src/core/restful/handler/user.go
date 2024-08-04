package handler

import (
	"context"

	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

type UserRestful struct {
	userService service.User
}

func NewUserRestful(us service.User) *UserRestful {
	return &UserRestful{
		userService: us,
	}
}

func (u *UserRestful) GetCurrent(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	res, err := u.userService.FindByEmail(context.Background(), email)
	if err != nil {
		return err
	}

	user := new(dto.SanitizedUserRes)
	copier.Copy(user, res)

	return c.Status(200).JSON(fiber.Map{"data": user})
}
