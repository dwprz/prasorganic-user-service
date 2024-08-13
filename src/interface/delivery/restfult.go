package delivery

import (
	"context"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type ImageKitRestful interface {
	UploadImage(ctx context.Context, path string, filename string) (*uploader.UploadResult, error)
	DeleteFile(ctx context.Context, fileId string)
}
