package server

import (
	"context"

	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserGrpcImpl struct {
	logger      *logrus.Logger
	userService service.User
	pb.UnimplementedUserServiceServer
}

func NewUserGrpc(l *logrus.Logger, us service.User) pb.UserServiceServer {
	return &UserGrpcImpl{
		logger:      l,
		userService: us,
	}
}

func (u *UserGrpcImpl) Create(ctx context.Context, ur *pb.RegisterRequest) (*emptypb.Empty, error) {
	data := &dto.CreateUserRequest{}
	if err := copier.Copy(data, ur); err != nil {
		return nil, err
	}

	if err := u.userService.Create(ctx, data); err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UserGrpcImpl) FindByEmail(ctx context.Context, e *pb.Email) (*pb.FindUserResponse, error) {
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

	return &pb.FindUserResponse{Data: user}, nil
}

func (u *UserGrpcImpl) Upsert(ctx context.Context, data *pb.LoginWithGoogleRequest) (*pb.User, error) {
	req := new(dto.UpsertUserRequest)
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
