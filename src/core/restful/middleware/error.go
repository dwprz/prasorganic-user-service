package middleware

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/common/errors/restful"
	"github.com/dwprz/prasorganic-user-service/src/common/helper"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/imagekit"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"google.golang.org/grpc/status"
)

func (m *Middleware) Error(c *fiber.Ctx, err error) error {
	restful.LogError(c, err)

	if err != nil && c.OriginalURL() == "/api/users/current/photo-profile" && c.Method() == "PATCH" {
		filename := c.Locals("filename").(string)
		if filename != "" {
			go helper.DeleteFile("./tmp/" + filename)
		}

		req, ok := c.Locals("upload_imagekit_result").(*uploader.UploadResult)
		if ok && req.FileId != "" {
			go imagekit.IK.Media.DeleteFile(context.Background(), req.FileId)
		}
	}

	if st, ok := status.FromError(err); ok {
		return restful.HandleGrpcError(c, st)
	}

	if validationError, ok := err.(validator.ValidationErrors); ok {
		return restful.HandleValidationError(c, validationError)
	}

	if responseError, ok := err.(*errors.Response); ok {
		return restful.HandleResponseError(c, responseError)
	}

	if jwtError := restful.HanldeJwtError(err); jwtError != nil {
		return c.Status(401).JSON(fiber.Map{"errors": jwtError.Error()})
	}

	if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		return restful.HandleJsonError(c, jsonError)
	}

	if strconvError, ok := err.(*strconv.NumError); ok {
		return restful.HandleStrconvError(c, strconvError)
	}

	return c.Status(500).JSON(fiber.Map{
		"errors": "sorry, internal server error try again later",
	})
}
