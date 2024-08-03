package repository

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

func (u *UserMock) FindByFields(ctx context.Context, fields *entity.User) (*entity.User, error) {
	arguments := u.Mock.Called(ctx, fields)

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

func (u *UserMock) UpdateRefreshToken(ctx context.Context, data *dto.UpdateRefreshToken) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}