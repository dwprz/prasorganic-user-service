package handler

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/dwprz/prasorganic-user-service/src/interface/helper"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

type UserRestful struct {
	userService service.User
	helper      helper.Helper
}

func NewUserRestful(us service.User, h helper.Helper) *UserRestful {
	return &UserRestful{
		userService: us,
		helper:      h,
	}
}

func (u *UserRestful) GetCurrent(c *fiber.Ctx) error {
	defer u.helper.HandlePanic(c)

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

func (u *UserRestful) UpdateProfile(c *fiber.Ctx) error {
	defer u.helper.HandlePanic(c)

	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := new(dto.UpdateProfileReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.Email = email

	res, err := u.userService.UpdateProfile(context.Background(), req)

	if err != nil {
		return err
	}

	user := new(dto.SanitizedUserRes)
	copier.Copy(user, res)

	return c.Status(200).JSON(fiber.Map{"data": user})
}

func (u *UserRestful) UpdatePassword(c *fiber.Ctx) error {
	defer u.helper.HandlePanic(c)

	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := new(dto.UpdatePasswordReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.Email = email
	err := u.userService.UpdatePassword(context.Background(), req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": "successfully updated the password"})
}

func (u *UserRestful) UpdateEmail(c *fiber.Ctx) error {
	defer u.helper.HandlePanic(c)

	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := new(dto.UpdateEmailReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.Email = email
	res, err := u.userService.UpdateEmail(context.Background(), req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "update_email",
		Value:    base64.StdEncoding.EncodeToString([]byte(res)),
		HTTPOnly: true,
		Path:     "/api/users/current/email/verify",
		Expires:  time.Now().Add(10 * time.Minute),
	})

	return c.Status(200).JSON(fiber.Map{"data": "Successfully requested email update"})
}

func (u *UserRestful) VerifyUpdateEmail(c *fiber.Ctx) error {
	defer u.helper.HandlePanic(c)

	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := new(dto.VerifyUpdateEmailReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	newEmail, err := base64.StdEncoding.DecodeString(c.Cookies("update_email"))
	if err != nil {
		return err
	}

	req.NewEmail = string(newEmail)
	req.Email = email

	res, err := u.userService.VerifyUpdateEmail(context.Background(), req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
	})

	c.Cookie(u.helper.ClearCookie("update_email", "/api/users/current/email/verify")) // clear cookie

	return c.Status(200).JSON(fiber.Map{"data": res.Data})
}

func (u *UserRestful) UpdatePhotoProfile(c *fiber.Ctx) error {
	defer u.helper.HandlePanic(c)

	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := c.Locals("update_photo_profile_req").(dto.UpdatePhotoProfileReq)
	req.Email = email

	res, err := u.userService.UpdatePhotoProfile(context.Background(), &req)
	if err != nil {
		return err
	}

	user := new(dto.SanitizedUserRes)
	copier.Copy(user, res)

	return c.Status(200).JSON(fiber.Map{"data": user})
}
