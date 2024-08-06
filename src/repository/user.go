package repository

import (
	"context"

	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/interface/cache"
	"github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// *caching menggunakan ctx.Background() supaya tidak dicancel, karena ada case context lintas server

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

func (u *UserImpl) Create(ctx context.Context, data *dto.CreateReq) error {
	query := "INSERT INTO users (user_id, email, full_name, password) VALUES($1, $2, $3, $4) RETURNING *;"

	user := new(entity.User)
	if err := u.db.WithContext(ctx).Raw(query, data.UserId, data.Email, data.FullName, data.Password).Scan(user).Error; err != nil {

		if errPG, ok := err.(*pgconn.PgError); ok && errPG.Code == "23505" {
			return &errors.Response{
				HttpCode: 409,
				GrpcCode: codes.AlreadyExists,
				Message:  "user already exists",
			}
		}

		return err
	}

	go u.userCache.Cache(context.Background(), user)

	return nil
}

func (u *UserImpl) FindByFields(ctx context.Context, fields *entity.User) (*entity.User, error) {
	user := new(entity.User)

	res := u.db.WithContext(ctx).Where(fields).Find(user)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	go u.userCache.Cache(context.Background(), user)

	return user, nil
}

func (u *UserImpl) Upsert(ctx context.Context, data *dto.UpsertReq) (*entity.User, error) {
	user := &entity.User{}

	query := `
	INSERT INTO 
		users (user_id, email, full_name, photo_profile, refresh_token, role, created_at)
	VALUES
		($1, $2, $3, $4, $5, 'USER', now())
	ON CONFLICT
		(email)
	DO UPDATE SET
		full_name = $3, updated_at = now()
	RETURNING *;
	`

	if err := u.db.WithContext(ctx).Raw(
		query,
		data.UserId,
		data.Email,
		data.FullName,
		data.PhotoProfile,
		data.RefreshToken).Scan(user).Error; err != nil {
		return nil, err
	}

	go u.userCache.Cache(context.Background(), user)

	return user, nil
}

func (u *UserImpl) UpdateByEmail(ctx context.Context, data *entity.User) (*entity.User, error) {
	user := new(entity.User)

	res := u.db.WithContext(ctx).Where("email = ?", data.Email).Updates(data).Clauses(clause.Returning{
		Columns: []clause.Column{{
			Name: "*",
		}},
	}).Find(user)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected > 0 {
		go u.userCache.Cache(context.Background(), user)
	}

	return user, nil
}

func (u *UserImpl) UpdateEmail(ctx context.Context, email string, newEmail string) (*entity.User, error) {
	user := new(entity.User)

	query := `UPDATE users SET email = $1, updated_at = now() WHERE email = $2 RETURNING *;`

	res := u.db.WithContext(ctx).Raw(query, newEmail, email).Scan(user)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected > 0 {
		go u.userCache.DeleteByEmail(context.Background(), email)
		go u.userCache.Cache(context.Background(), user)
	}

	return user, nil
}

func (u *UserImpl) SetNullRefreshToken(ctx context.Context, refreshToken string) error {
	user := new(entity.User)

	query := `UPDATE users SET refresh_token = NULL, updated_at = now() WHERE refresh_token = $1 RETURNING *;`

	res := u.db.WithContext(ctx).Raw(query, refreshToken).Scan(user)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected > 0 {
		go u.userCache.Cache(context.Background(), user)
	}

	return nil
}
