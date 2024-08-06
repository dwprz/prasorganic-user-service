package service

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/stretchr/testify/mock"
)

type UserMock struct {
	mock.Mock
}

func NewUserMock() *UserMock {
	return &UserMock{
		Mock: mock.Mock{},
	}
}

func (u *UserMock) Create(ctx context.Context, data *dto.CreateReq) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (u *UserMock) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	arguments := u.Mock.Called(ctx, email)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.User), arguments.Error(1)
}

func (u *UserMock) FindByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error) {
	arguments := u.Mock.Called(ctx, refreshToken)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.User), arguments.Error(1)
}

func (u *UserMock) Upsert(ctx context.Context, data *dto.UpsertReq) (*entity.User, error) {
	arguments := u.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.User), arguments.Error(1)
}

func (u *UserMock) UpdateProfile(ctx context.Context, data *dto.UpdateProfileReq) (*entity.User, error) {
	arguments := u.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.User), arguments.Error(1)
}

func ( u *UserMock) UpdatePhotoProfile(ctx context.Context, data *dto.UpdatePhotoProfileReq) (*entity.User, error) {
	arguments := u.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.User), arguments.Error(1)
}

func (u *UserMock) UpdatePassword(ctx context.Context, data *dto.UpdatePasswordReq) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (u *UserMock) UpdateEmail(ctx context.Context, data *dto.UpdateEmailReq) (newEmail string, err error) {
	arguments := u.Mock.Called(ctx, data)

	return arguments.String(0), arguments.Error(1)
}

func (u *UserMock) VerifyUpdateEmail(ctx context.Context, data *dto.VerifyUpdateEmailReq) (*dto.VerifyUpdateEmailRes, error) {
	arguments := u.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*dto.VerifyUpdateEmailRes), arguments.Error(1)
}

func (u *UserMock) AddRefreshToken(ctx context.Context, data *dto.AddRefreshTokenReq) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (u *UserMock) SetNullRefreshToken(ctx context.Context, refreshToken string) error {
	arguments := u.Mock.Called(ctx, refreshToken)

	return arguments.Error(0)
}
