package repository

import (
	"context"

	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
)

type User interface {
	Create(ctx context.Context, data *dto.CreateReq) error
	FindByFields(ctx context.Context, fields *entity.User) (*entity.User, error)
	Upsert(ctx context.Context, data *dto.UpsertReq) (*entity.User, error)
	AddRefreshToken(ctx context.Context, data *dto.AddRefreshTokenReq) error
	SetNullRefreshToken(ctx context.Context, refreshToken string) error
}
