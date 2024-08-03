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

func (s *UserImpl) Create(ctx context.Context, data *dto.CreateUserRequest) error {
	if err := s.validate.Struct(data); err != nil {
		return err
	}

	if err := s.userRepository.Create(ctx, data); err != nil {
		return err
	}

	return nil
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

func (s *UserImpl) Upsert(ctx context.Context, data *dto.UpsertUserRequest) (*entity.User, error) {
	if err := s.validate.Struct(data); err != nil {
		return nil, err
	}

	user, err := s.userRepository.Upsert(ctx, data)
	if err != nil {
		return nil, err
	}

	return user, nil
}
