package repository

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
)

type User interface {
	Create(ctx context.Context, data *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
