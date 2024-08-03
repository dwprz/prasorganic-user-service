package service

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
)

type User interface {
	Create(ctx context.Context, data *dto.CreateUserRequest) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Upsert(ctx context.Context, data *dto.UpsertUserRequest) (*entity.User, error)
}
