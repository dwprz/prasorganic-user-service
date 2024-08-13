package test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/mock/service"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// go test -v ./src/core/grpc/handler/test/... -count=1 -p=1
// go test -run ^TestServer_FindUserByEmail$ -v ./src/core/grpc/handler/test -count=1

type FindUserByEmailTestSuite struct {
	suite.Suite
	grpcServer       *server.Grpc
	userGrpcDelivery pb.UserServiceClient
	userGrpcConn     *grpc.ClientConn
	userService      *service.UserMock
}

func (f *FindUserByEmailTestSuite) SetupSuite() {
	f.userService = service.NewUserMock()
	f.grpcServer = util.InitGrpcServerTest(f.userService)

	go f.grpcServer.Run()

	time.Sleep(1 * time.Second)

	userGrpcDelivery, userGrpcConn := util.InitUserGrpcDelivery()
	f.userGrpcDelivery = userGrpcDelivery
	f.userGrpcConn = userGrpcConn
}

func (f *FindUserByEmailTestSuite) TearDownSuite() {
	f.grpcServer.Stop()
	f.userGrpcConn.Close()
}

func (f *FindUserByEmailTestSuite) Test_Success() {
	request := &pb.Email{Email: "johndoe@gmail.com"}

	user := &entity.User{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
	}

	f.userService.Mock.On("FindByEmail", mock.Anything, request.Email).Return(user, nil)

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Basic "+auth)

	res, err := f.userGrpcDelivery.FindByEmail(ctx, request)
	assert.NoError(f.T(), err)
	assert.NotNil(f.T(), res.Data)
}

func (f *FindUserByEmailTestSuite) Test_NotFound() {
	request := &pb.Email{Email: "notfounduser@gmail.com"}

	f.userService.Mock.On("FindByEmail", mock.Anything, request.Email).Return(nil, nil)

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Basic "+auth)

	res, err := f.userGrpcDelivery.FindByEmail(ctx, request)
	assert.NoError(f.T(), err)
	assert.Nil(f.T(), res.Data)
}

func TestServer_FindUserByEmail(t *testing.T) {
	suite.Run(t, new(FindUserByEmailTestSuite))
}
