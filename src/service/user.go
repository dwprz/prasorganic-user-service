package service

import (
	"context"

	"github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/common/helper"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/client"
	v "github.com/dwprz/prasorganic-user-service/src/infrastructure/validator"
	"github.com/dwprz/prasorganic-user-service/src/interface/cache"
	"github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

type UserImpl struct {
	grpcClient     *client.Grpc
	userRepository repository.User
	userCache      cache.User
}

func NewUser(gc *client.Grpc, ur repository.User, uc cache.User) service.User {
	return &UserImpl{
		grpcClient:     gc,
		userRepository: ur,
		userCache:      uc,
	}
}

func (u *UserImpl) Create(ctx context.Context, data *dto.CreateReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	err := u.userRepository.Create(ctx, data)
	return err
}

func (u *UserImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if err := v.Validate.VarCtx(ctx, email, `required,email,min=10,max=100`); err != nil {
		return nil, err
	}

	if userCache := u.userCache.FindByEmail(ctx, email); userCache != nil {
		return userCache, nil
	}

	res, err := u.userRepository.FindByFields(ctx, &entity.User{Email: email})
	return res, err
}

func (u *UserImpl) FindByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error) {
	if err := v.Validate.VarCtx(ctx, refreshToken, `required,min=50,max=500`); err != nil {
		return nil, err
	}

	res, err := u.userRepository.FindByFields(ctx, &entity.User{
		RefreshToken: refreshToken,
	})

	return res, err
}

func (u *UserImpl) Upsert(ctx context.Context, data *dto.UpsertReq) (*entity.User, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	res, err := u.userRepository.Upsert(ctx, data)
	return res, err
}

func (u *UserImpl) UpdateProfile(ctx context.Context, data *dto.UpdateProfileReq) (*entity.User, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	user := u.userCache.FindByEmail(ctx, data.Email)
	if user == nil {
		res, err := u.userRepository.FindByFields(ctx, &entity.User{
			Email: data.Email,
		})

		if err != nil {
			return nil, err
		}

		if res == nil {
			return nil, &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "user not found"}
		}

		user = res
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return nil, &errors.Response{HttpCode: 400, GrpcCode: codes.InvalidArgument, Message: "password is invalid"}
	}

	res, err := u.userRepository.UpdateByEmail(ctx, &entity.User{
		Email:    data.Email,
		FullName: data.FullName,
		Whatsapp: data.Whatsapp,
	})

	return res, err
}

func (u *UserImpl) UpdatePassword(ctx context.Context, data *dto.UpdatePasswordReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	user := u.userCache.FindByEmail(ctx, data.Email)
	if user == nil {
		res, err := u.userRepository.FindByFields(ctx, &entity.User{
			Email: data.Email,
		})

		if err != nil {
			return err
		}

		if res == nil {
			return &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "user not found"}
		}

		user = res
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return &errors.Response{HttpCode: 400, GrpcCode: codes.InvalidArgument, Message: "password is invalid"}
	}

	encryptPwd, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = u.userRepository.UpdateByEmail(ctx, &entity.User{
		Email:    data.Email,
		Password: string(encryptPwd),
	})

	return err
}

func (u *UserImpl) UpdateEmail(ctx context.Context, data *dto.UpdateEmailReq) (newEmail string, err error) {
	if err := v.Validate.Struct(data); err != nil {
		return "", err
	}

	user := u.userCache.FindByEmail(ctx, data.Email)
	if user == nil {
		res, err := u.userRepository.FindByFields(ctx, &entity.User{
			Email: data.Email,
		})

		if err != nil {
			return "", err
		}

		if res == nil {
			return "", &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "user not found"}
		}

		user = res
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return "", &errors.Response{HttpCode: 400, GrpcCode: codes.InvalidArgument, Message: "password is invalid"}
	}

	go u.grpcClient.Otp.Send(ctx, data.NewEmail)

	return data.NewEmail, nil
}

func (u *UserImpl) VerifyUpdateEmail(ctx context.Context, data *dto.VerifyUpdateEmailReq) (*dto.VerifyUpdateEmailRes, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	verifyRes, err := u.grpcClient.Otp.Verify(ctx, &otp.VerifyReq{
		Email: data.NewEmail,
		Otp:   data.Otp,
	})

	if err != nil {
		return nil, err
	}

	if !verifyRes.Valid {
		return nil, &errors.Response{HttpCode: 400, Message: "otp is invalid"}
	}

	res, err := u.userRepository.UpdateEmail(ctx, data.Email, data.NewEmail)
	if err != nil {
		return nil, err
	}

	accessToken, err := helper.GenerateAccessToken(res.UserId, res.Email, res.Role)
	if err != nil {
		return nil, err
	}

	user := new(dto.SanitizedUserRes)
	if err := copier.Copy(user, res); err != nil {
		return nil, err
	}

	return &dto.VerifyUpdateEmailRes{
		Data:        user,
		AccessToken: accessToken,
	}, nil
}

func (u *UserImpl) UpdatePhotoProfile(ctx context.Context, data *dto.UpdatePhotoProfileReq) (*entity.User, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	res, err := u.userRepository.UpdateByEmail(ctx, &entity.User{
		Email:          data.Email,
		PhotoProfileId: data.PhotoProfileId,
		PhotoProfile:   data.PhotoProfile,
	})

	return res, err
}

func (u *UserImpl) AddRefreshToken(ctx context.Context, data *dto.AddRefreshTokenReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	_, err := u.userRepository.UpdateByEmail(ctx, &entity.User{
		Email:        data.Email,
		RefreshToken: data.RefreshToken,
	})

	return err
}

func (u *UserImpl) SetNullRefreshToken(ctx context.Context, refreshToken string) error {
	if err := v.Validate.VarCtx(ctx, refreshToken, `required,min=50,max=500`); err != nil {
		return err
	}

	err := u.userRepository.SetNullRefreshToken(ctx, refreshToken)
	return err
}
