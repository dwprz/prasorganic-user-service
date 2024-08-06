package service

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
)

type User interface {
	Create(ctx context.Context, data *dto.CreateReq) error
	
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error)
	
	Upsert(ctx context.Context, data *dto.UpsertReq) (*entity.User, error)
	UpdateProfile(ctx context.Context, data *dto.UpdateProfileReq) (*entity.User, error)
	UpdatePhotoProfile(ctx context.Context, data *dto.UpdatePhotoProfileReq) (*entity.User, error)
	UpdatePassword(ctx context.Context, data *dto.UpdatePasswordReq) error
	UpdateEmail(ctx context.Context, data *dto.UpdateEmailReq) (newEmail string, err error)
	
	VerifyUpdateEmail(ctx context.Context, data *dto.VerifyUpdateEmailReq) (*dto.VerifyUpdateEmailRes, error)
	AddRefreshToken(ctx context.Context, data *dto.AddRefreshTokenReq) error
	SetNullRefreshToken(ctx context.Context, refreshToken string) error
}
