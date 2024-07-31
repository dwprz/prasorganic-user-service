package service

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/interface/cache"
	"github.com/dwprz/prasorganic-user-service/interface/repository"
	"github.com/dwprz/prasorganic-user-service/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
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

func (s *UserImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if userCache := s.userCache.FindByEmail(ctx, email); userCache != nil {
		return userCache, nil
	}

	result, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserImpl) Create(ctx context.Context, data *dto.UserCreate) error {
	if err := s.validate.Struct(data); err != nil {
		return err
	}

	userEntity := new(entity.User)
	if err := copier.Copy(userEntity, data); err != nil {
		return err
	}

	if err := s.userRepository.Create(ctx, userEntity); err != nil {
		return err
	}

	return nil
}
