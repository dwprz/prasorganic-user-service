package test

import (
	"context"
	"testing"
	"time"

	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/common/helper"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	grpcapp "github.com/dwprz/prasorganic-user-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/imagekit"
	"github.com/dwprz/prasorganic-user-service/src/mock/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// go test -v ./src/core/grpc/server/test/... -count=1 -p=1
// go test -run ^TestServer_CreateUser$ -v ./src/core/grpc/server/test -count=1

type CreateUserTestSuite struct {
	suite.Suite
	grpcServer     *grpcapp.Server
	userGrpcClient pb.UserServiceClient
	userGrpcConn   *grpc.ClientConn
	userService    *service.UserMock
	logger         *logrus.Logger
}

func (c *CreateUserTestSuite) SetupSuite() {
	c.logger = logger.New()
	conf := config.New("DEVELOPMENT", c.logger)
	imageKit := imagekit.New(conf)
	helper := helper.New(imageKit, conf, c.logger)

	// mock
	c.userService = service.NewUserMock()

	userGrpcServer := server.NewUserGrpc(c.logger, c.userService)
	unaryResponseInterceptor := interceptor.NewUnaryResponse(c.logger, helper)

	c.grpcServer = grpcapp.NewServer(conf.CurrentApp.GrpcPort, userGrpcServer, unaryResponseInterceptor, c.logger)

	go c.grpcServer.Run()

	time.Sleep(2 * time.Second)

	grpcAddress := "localhost:" + conf.CurrentApp.GrpcPort
	userGrpcClient, userGrpcConn := util.NewGrpcUserClient(grpcAddress)

	c.userGrpcClient = userGrpcClient
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

	registerReq := new(pb.RegisterRequest)
	err := copier.Copy(registerReq, userCreate)
	assert.NoError(c.T(), err)

	_, err = c.userGrpcClient.Create(context.Background(), registerReq)
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

	registerReq := new(pb.RegisterRequest)
	err := copier.Copy(registerReq, userCreate)
	assert.NoError(c.T(), err)

	_, err = c.userGrpcClient.Create(context.Background(), registerReq)
	assert.Error(c.T(), err)
}

func TestServer_CreateUser(t *testing.T) {
	suite.Run(t, new(CreateUserTestSuite))
}
