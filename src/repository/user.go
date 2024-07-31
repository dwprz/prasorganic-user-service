package repository

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/interface/cache"
	"github.com/dwprz/prasorganic-user-service/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type UserImpl struct {
	db        *gorm.DB
	userCache cache.User
}

func NewUser(db *gorm.DB, uc cache.User) repository.User {
	return &UserImpl{
		db:        db,
		userCache: uc,
	}
}

func (u *UserImpl) Create(ctx context.Context, data *entity.User) error {
	query := "INSERT INTO users (email, full_name, password) VALUES($1, $2, $3) RETURNING *;"

	user := new(entity.User)
	if err := u.db.WithContext(ctx).Raw(query, data.Email, data.FullName, data.Password).Scan(user).Error; err != nil {

		if errPG, ok := err.(*pgconn.PgError); ok && errPG.Code == "23505" {
			return &errors.Response{
				HttpCode: 409,
				GrpcCode: codes.AlreadyExists,
				Message:  "user already exists",
			}
		}

		return err
	}

	u.userCache.Cache(ctx, user)

	return nil
}

func (u *UserImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}

	result := u.db.WithContext(ctx).Raw("SELECT * FROM users WHERE email = $1;", email).Scan(user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	u.userCache.Cache(ctx, user)

	return user, nil
}