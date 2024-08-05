package service

import (
	"context"
	"github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-user-service/src/interface/cache"
	"github.com/dwprz/prasorganic-user-service/src/interface/helper"
	"github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

type UserImpl struct {
	grpcClient     *grpc.Client
	validate       *validator.Validate
	userRepository repository.User
	userCache      cache.User
	helper         helper.Helper
}

func NewUser(gc *grpc.Client, v *validator.Validate, ur repository.User, uc cache.User, h helper.Helper) service.User {
	return &UserImpl{
		grpcClient:     gc,
		validate:       v,
		userRepository: ur,
		userCache:      uc,
		helper:         h,
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

func (u *UserImpl) UpdateProfile(ctx context.Context, data *dto.UpdateProfileReq) (*entity.User, error) {
	if err := u.validate.Struct(data); err != nil {
		return nil, err
	}

	res, err := u.userRepository.FindByFields(ctx, &entity.User{
		Email: data.Email,
	})

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "user not found"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(data.Password)); err != nil {
		return nil, &errors.Response{HttpCode: 400, GrpcCode: codes.InvalidArgument, Message: "password is invalid"}
	}

	res, err = u.userRepository.UpdateByEmail(ctx, &entity.User{
		Email:    data.Email,
		FullName: data.FullName,
		Whatsapp: data.Whatsapp,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UserImpl) UpdatePassword(ctx context.Context, data *dto.UpdatePasswordReq) error {
	if err := u.validate.Struct(data); err != nil {
		return err
	}

	res, err := u.userRepository.FindByFields(ctx, &entity.User{
		Email: data.Email,
	})

	if err != nil {
		return err
	}

	if res == nil {
		return &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "user not found"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(data.Password)); err != nil {
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

	if err != nil {
		return err
	}

	return nil
}

func (u *UserImpl) UpdateEmail(ctx context.Context, data *dto.UpdateEmailReq) (newEmail string, err error) {
	if err := u.validate.Struct(data); err != nil {
		return "", err
	}

	res, err := u.userRepository.FindByFields(ctx, &entity.User{
		Email: data.Email,
	})

	if err != nil {
		return "", err
	}

	if res == nil {
		return "", &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "user not found"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(data.Password)); err != nil {
		return "", &errors.Response{HttpCode: 400, GrpcCode: codes.InvalidArgument, Message: "password is invalid"}
	}

	go u.grpcClient.Otp.Send(ctx, data.NewEmail)

	return data.NewEmail, nil
}

func (u *UserImpl) VerifyUpdateEmail(ctx context.Context, data *dto.VerifyUpdateEmailReq) (*dto.VerifyUpdateEmailRes, error) {
	if err := u.validate.Struct(data); err != nil {
		return nil, err
	}

	verifyRes, err := u.grpcClient.Otp.Verify(ctx, &otp.VerifyRequest{
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

	accessToken, err := u.helper.GenerateAccessToken(res.UserId, res.Email, res.Role)
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

func (u *UserImpl) AddRefreshToken(ctx context.Context, data *dto.AddRefreshTokenReq) error {
	if err := u.validate.Struct(data); err != nil {
		return err
	}

	_, err := u.userRepository.UpdateByEmail(ctx, &entity.User{
		Email:        data.Email,
		RefreshToken: data.RefreshToken,
	})

	if err != nil {
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
