package delivery

import (
	"context"

	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"github.com/stretchr/testify/mock"
)

type ImageKitRestfulMock struct {
	mock.Mock
}

func NewImageKitRestfulMock() *ImageKitRestfulMock {
	return &ImageKitRestfulMock{
		Mock: mock.Mock{},
	}
}

func (i *ImageKitRestfulMock) UploadImage(ctx context.Context, path string, filename string) (*uploader.UploadResult, error) {
	arguments := i.Mock.Called(ctx, path, filename)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*uploader.UploadResult), arguments.Error(1)
}

func DeleteFile(ctx context.Context, fileId string) {}
