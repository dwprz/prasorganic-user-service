package delivery

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/imagekit"
	"github.com/dwprz/prasorganic-user-service/src/interface/delivery"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"github.com/sirupsen/logrus"
)

type ImageKit struct {}

func NewImageKit() delivery.ImageKitRestful {
	return &ImageKit{}
}

func (i *ImageKit) UploadImage(ctx context.Context, path string, filename string) (*uploader.UploadResult, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	base64String := base64.StdEncoding.EncodeToString(fileData)
	file := "data:image/jpeg;base64," + base64String

	useUniqueFileName := false

	res, err := imagekit.IK.Uploader.Upload(ctx, file, uploader.UploadParam{
		FileName:          filename,
		UseUniqueFileName: &useUniqueFileName,
	})

	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (i *ImageKit) DeleteFile(ctx context.Context, fileId string) {
	_, err := imagekit.IK.Media.DeleteFile(ctx, fileId)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.ImageKit/DeleteFile", "section": "ik.Media.DeleteFile"}).Error(err)
	}
}
