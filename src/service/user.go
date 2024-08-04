package service

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/src/interface/cache"
	"github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/go-playground/validator/v10"
)

type UserImpl struct {
	validate       *validator.Validate
	userRepository repository.User
	userCache      cache.User
}

func NewUser(v *validator.Validate, ur repository.User, uc cache.User) service.User {
	return &UserImpl{
		validate:       v,
		userRepository: ur,
		userCache:      uc,
	}
}

func (u *UserImpl) Create(ctx context.Context, data *dto.CreateReq) error {
	if err := u.validate.Struct(data); err != nil {
		return err
	}

	if err := u.userRepository.Create(ctx, data); err != nil {
		return err
	}

	return nil
}

func (u *UserImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if err := u.validate.VarCtx(ctx, email, `required,email,min=10,max=100`); err != nil {
		return nil, err
	}

	if userCache := u.userCache.FindByEmail(ctx, email); userCache != nil {
		return userCache, nil
	}

	res, err := u.userRepository.FindByFields(ctx, &entity.User{Email: email})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UserImpl) FindByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error) {
	if err := u.validate.VarCtx(ctx, refreshToken, `required,min=50,max=500`); err != nil {
		return nil, err
	}

	res, err := u.userRepository.FindByFields(ctx, &entity.User{
		RefreshToken: refreshToken,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UserImpl) Upsert(ctx context.Context, data *dto.UpsertReq) (*entity.User, error) {
	if err := u.validate.Struct(data); err != nil {
		return nil, err
	}

	res, err := u.userRepository.Upsert(ctx, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UserImpl) AddRefreshToken(ctx context.Context, data *dto.AddRefreshTokenReq) error {
	if err := u.validate.Struct(data); err != nil {
		return err
	}

	if err := u.userRepository.AddRefreshToken(ctx, data); err != nil {
		return err
	}

	return nil
}

func (u *UserImpl) SetNullRefreshToken(ctx context.Context, refreshToken string) error {
	if err := u.validate.VarCtx(ctx, refreshToken, `required,min=50,max=500`); err != nil {
		return err
	}

	if err := u.userRepository.SetNullRefreshToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}
