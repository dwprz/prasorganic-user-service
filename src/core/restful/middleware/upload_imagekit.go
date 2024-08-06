package middleware

import (
	"encoding/base64"
	"os"

	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

func (m *Middleware) UploadToImageKit(c *fiber.Ctx) error {
	photoProfileId := c.FormValue("photo_profile_id")
	if photoProfileId != "" {
		go m.imageKit.Media.DeleteFile(c.Context(), photoProfileId)
	}

	filename := c.Locals("filename").(string)

	fileData, err := os.ReadFile("./tmp/" + filename)
	if err != nil {
		return err
	}

	base64String := base64.StdEncoding.EncodeToString(fileData)
	file := "data:image/jpeg;base64," + base64String

	useUniqueFileName := false

	res, err := m.imageKit.Uploader.Upload(c.Context(), file, uploader.UploadParam{
		FileName:          filename,
		UseUniqueFileName: &useUniqueFileName,
	})

	if err != nil {
		return err
	}

	req := dto.UpdatePhotoProfileReq{
		NewPhotoProfileId: res.Data.FileId,
		NewPhotoProfile:   res.Data.Url,
	}

	c.Locals("update_photo_profile_req", req)
	go m.helper.DeleteFile("./tmp/" + filename)

	return c.Next()
}
