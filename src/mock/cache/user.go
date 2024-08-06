package cache

import (
	"context"
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

func (u *UserMock) Cache(ctx context.Context, user *entity.User) {}

func (u *UserMock) FindByEmail(ctx context.Context, email string) *entity.User {
	arguments := u.Mock.Called(ctx, email)

	if arguments.Get(0) == nil {
		return nil
	}

	return arguments.Get(0).(*entity.User)
}

func (u *UserMock) DeleteByEmail(ctx context.Context, email string) {}