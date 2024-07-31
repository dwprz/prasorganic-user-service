package server

import (
	"context"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/interface/service"
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

func (g *UserGrpcImpl) FindByEmail(ctx context.Context, u *pb.Email) (*pb.FindUserResponse, error) {
	result, err := g.userService.FindByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	user := new(pb.User)
	if err := copier.Copy(user, result); err != nil {
		return nil, err
	}

	user.CreatedAt = timestamppb.New(result.CreatedAt)
	user.UpdatedAt = timestamppb.New(result.UpdatedAt)

	return &pb.FindUserResponse{Data: user}, nil
}

func (g *UserGrpcImpl) Create(ctx context.Context, ur *pb.RegisterRequest) (*emptypb.Empty, error) {
	data := &dto.UserCreate{}
	if err := copier.Copy(data, ur); err != nil {
		return nil, err
	}

	if err := g.userService.Create(ctx, data); err != nil {
		return nil, err
	}

	return nil, nil
}
