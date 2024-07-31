package service

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
)

type User interface {
	Create(ctx context.Context, data *dto.UserCreate) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}