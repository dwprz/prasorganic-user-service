package test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/mock/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// go test -v ./src/core/grpc/handler/test/... -count=1 -p=1
// go test -run ^TestServer_CreateUser$ -v ./src/core/grpc/handler/test -count=1

type CreateUserTestSuite struct {
	suite.Suite
	grpcServer       *server.Grpc
	userGrpcDelivery pb.UserServiceClient
	userGrpcConn     *grpc.ClientConn
	userService      *service.UserMock
}

func (c *CreateUserTestSuite) SetupSuite() {
	c.userService = service.NewUserMock()
	c.grpcServer = util.InitGrpcServerTest(c.userService)

	go c.grpcServer.Run()

	time.Sleep(1 * time.Second)

	userGrpcDelivery, userGrpcConn := util.InitUserGrpcDelivery()
	c.userGrpcDelivery = userGrpcDelivery
	c.userGrpcConn = userGrpcConn
}

func (c *CreateUserTestSuite) TearDownSuite() {
	c.grpcServer.Stop()
	c.userGrpcConn.Close()
}

func (c *CreateUserTestSuite) Test_Success() {
	userCreate := &dto.CreateReq{
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	c.userService.Mock.On("Create", mock.Anything, userCreate).Return(nil)

	registerReq := new(pb.RegisterReq)
	err := copier.Copy(registerReq, userCreate)
	assert.NoError(c.T(), err)

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Basic "+auth)

	_, err = c.userGrpcDelivery.Create(ctx, registerReq)
	assert.NoError(c.T(), err)
}

func (c *CreateUserTestSuite) Test_AlreadyExists() {
	userCreate := &dto.CreateReq{
		Email:    "existeduser@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	errorRes := &errors.Response{HttpCode: 409, GrpcCode: codes.AlreadyExists}
	c.userService.Mock.On("Create", mock.Anything, userCreate).Return(errorRes)

	registerReq := new(pb.RegisterReq)
	err := copier.Copy(registerReq, userCreate)
	assert.NoError(c.T(), err)

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Basic "+auth)

	_, err = c.userGrpcDelivery.Create(ctx, registerReq)
	assert.Error(c.T(), err)

	st, ok :=status.FromError(err)
	assert.True(c.T(), ok)

	assert.Equal(c.T(), st.Code(), errorRes.GrpcCode)
}

func TestServer_CreateUser(t *testing.T) {
	suite.Run(t, new(CreateUserTestSuite))
}
