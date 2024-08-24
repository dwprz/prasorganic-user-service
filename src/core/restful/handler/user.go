package handler

import (
	"encoding/base64"
	"time"

	"github.com/dwprz/prasorganic-user-service/src/common/helper"
	"github.com/dwprz/prasorganic-user-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"github.com/jinzhu/copier"
)

type User struct {
	userService   service.User
	restfulClient *client.Restful
}

func NewUser(us service.User, rc *client.Restful) *User {
	return &User{
		userService:   us,
		restfulClient: rc,
	}
}

func (u *User) GetCurrent(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	res, err := u.userService.FindByEmail(c.Context(), email)
	if err != nil {
		return err
	}

	user := new(dto.SanitizedUserRes)
	copier.Copy(user, res)

	return c.Status(200).JSON(fiber.Map{"data": user})
}

func (u *User) UpdateProfile(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := new(dto.UpdateProfileReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.Email = email

	res, err := u.userService.UpdateProfile(c.Context(), req)

	if err != nil {
		return err
	}

	user := new(dto.SanitizedUserRes)
	if err := copier.Copy(user, res); err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": user})
}

func (u *User) UpdatePassword(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := new(dto.UpdatePasswordReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.Email = email
	err := u.userService.UpdatePassword(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": "successfully updated the password"})
}

func (u *User) UpdateEmail(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := new(dto.UpdateEmailReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.Email = email
	res, err := u.userService.UpdateEmail(c.Context(), req)
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

	return c.Status(200).JSON(fiber.Map{"data": "successfully requested email update"})
}

func (u *User) VerifyUpdateEmail(c *fiber.Ctx) error {
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

	res, err := u.userService.VerifyUpdateEmail(c.Context(), req)
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

	c.Cookie(helper.ClearCookie("update_email", "/api/users/current/email/verify")) // clear cookie

	return c.Status(200).JSON(fiber.Map{"data": res.Data})
}

func (u *User) UpdatePhotoProfile(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	email := userData["email"].(string)

	req := new(dto.UpdatePhotoProfileReq)

	uploadRes := c.Locals("upload_imagekit_result").(*uploader.UploadResult)
	req.PhotoProfileId = uploadRes.FileId
	req.PhotoProfile = uploadRes.Url
	req.Email = email

	res, err := u.userService.UpdatePhotoProfile(c.Context(), req)
	if err != nil {
		return err
	}

	user := new(dto.SanitizedUserRes)
	if err := copier.Copy(user, res); err != nil {
		return err
	}

	photoProfileId := c.FormValue("photo_profile_id")
	if photoProfileId != "" {
		go u.restfulClient.ImageKit.DeleteFile(c.Context(), photoProfileId)
	}

	return c.Status(200).JSON(fiber.Map{"data": user})
}
