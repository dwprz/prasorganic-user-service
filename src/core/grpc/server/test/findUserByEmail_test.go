package test

import (
	"context"
	"testing"
	"time"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/mock/service"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	grpcapp "github.com/dwprz/prasorganic-user-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -v ./src/core/grpc/server/test/... -count=1 -p=1
// go test -run ^TestServer_FindUserByEmail$ -v ./src/core/grpc/server/test -count=1

type FindUserByEmailTestSuite struct {
	suite.Suite
	grpcServer     *grpcapp.Server
	userGrpcClient user.UserServiceClient
	userGrpcConn   *grpc.ClientConn
	userService    *service.UserMock
	logger         *logrus.Logger
}

func (f *FindUserByEmailTestSuite) SetupSuite() {
	f.logger = logger.New()
	conf := config.New("DEVELOPMENT", f.logger)

	// mock
	f.userService = service.NewUserMock()

	userGrpcServer := server.NewUserGrpc(f.logger, f.userService)
	unaryResponseInterceptor := interceptor.NewUnaryResponse(f.logger)

	f.grpcServer = grpcapp.NewServer(conf.CurrentApp.GrpcPort, userGrpcServer, unaryResponseInterceptor, f.logger)

	go f.grpcServer.Run()

	time.Sleep(2 * time.Second)

	grpcAddress := "localhost:" + conf.CurrentApp.GrpcPort
	userGrpcClient, userGrpcConn := util.NewGrpcUserClient(grpcAddress)
	
	f.userGrpcClient = userGrpcClient
	f.userGrpcConn = userGrpcConn
}

func (f *FindUserByEmailTestSuite) TearDownSuite() {
	f.grpcServer.Stop()
	f.userGrpcConn.Close()
}

func (f *FindUserByEmailTestSuite) Test_Success() {
	request := &user.Email{Email: "johndoe@gmail.com"}

	user := &entity.User{
		UserID:   1,
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
	}

	f.userService.Mock.On("FindByEmail", mock.Anything, request.Email).Return(user, nil)
	
	res, err := f.userGrpcClient.FindByEmail(context.Background(), request)
	assert.NoError(f.T(), err)
	assert.NotNil(f.T(), res.Data)
}

func (f *FindUserByEmailTestSuite) Test_NotFound() {
	request := &user.Email{Email: "notfounduser@gmail.com"}

	f.userService.Mock.On("FindByEmail", mock.Anything, request.Email).Return(nil, nil)
	
	res, err := f.userGrpcClient.FindByEmail(context.Background(), request)
	assert.NoError(f.T(), err)
	assert.Nil(f.T(), res.Data)
}

func TestServer_FindUserByEmail(t *testing.T) {
	suite.Run(t, new(FindUserByEmailTestSuite))
}
