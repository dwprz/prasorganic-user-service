package handler

import (
	"context"

	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserGrpcImpl struct {
	userService service.User
	pb.UnimplementedUserServiceServer
}

func NewUserGrpc(us service.User) pb.UserServiceServer {
	return &UserGrpcImpl{
		userService: us,
	}
}

func (u *UserGrpcImpl) Create(ctx context.Context, ur *pb.RegisterReq) (*emptypb.Empty, error) {
	data := &dto.CreateReq{}
	if err := copier.Copy(data, ur); err != nil {
		return nil, err
	}

	if err := u.userService.Create(ctx, data); err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UserGrpcImpl) FindByEmail(ctx context.Context, e *pb.Email) (*pb.FindUserRes, error) {
	res, err := u.userService.FindByEmail(ctx, e.Email)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	user := new(pb.User)
	if err := copier.Copy(user, res); err != nil {
		return nil, err
	}

	user.CreatedAt = timestamppb.New(res.CreatedAt)
	user.UpdatedAt = timestamppb.New(res.UpdatedAt)

	return &pb.FindUserRes{Data: user}, nil
}

func (u *UserGrpcImpl) FindByRefreshToken(ctx context.Context, t *pb.RefreshToken) (*pb.FindUserRes, error) {
	res, err := u.userService.FindByRefreshToken(ctx, t.Token)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	user := new(pb.User)
	if err := copier.Copy(user, res); err != nil {
		return nil, err
	}

	user.CreatedAt = timestamppb.New(res.CreatedAt)
	user.UpdatedAt = timestamppb.New(res.UpdatedAt)

	return &pb.FindUserRes{Data: user}, nil
}

func (u *UserGrpcImpl) Upsert(ctx context.Context, data *pb.LoginWithGoogleReq) (*pb.User, error) {
	req := new(dto.UpsertReq)
	if err := copier.Copy(req, data); err != nil {
		return nil, err
	}

	res, err := u.userService.Upsert(ctx, req)
	if err != nil {
		return nil, err
	}

	user := new(pb.User)
	if err := copier.Copy(user, res); err != nil {
		return nil, err
	}

	user.CreatedAt = timestamppb.New(res.CreatedAt)
	user.UpdatedAt = timestamppb.New(res.UpdatedAt)

	return user, nil
}

func (u *UserGrpcImpl) AddRefreshToken(ctx context.Context, t *pb.AddRefreshTokenReq) (*emptypb.Empty, error) {
	req := &dto.AddRefreshTokenReq{
		Email:        t.Email,
		RefreshToken: t.Token,
	}

	err := u.userService.AddRefreshToken(ctx, req)
	return nil, err
}

func (u *UserGrpcImpl) SetNullRefreshToken(ctx context.Context, t *pb.RefreshToken) (*emptypb.Empty, error) {

	err := u.userService.SetNullRefreshToken(ctx, t.Token)
	return nil, err
}
