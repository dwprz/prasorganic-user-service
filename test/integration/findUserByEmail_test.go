package integration_test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"
	userpb "github.com/dwprz/prasorganic-proto/protogen/user"
	grpcapp "github.com/dwprz/prasorganic-user-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// go test -v ./tests/integration -count=1
// go test -run ^TestIntegration_FindUserByEmail$  -v ./test/integration -count=1

type FindUserByEmailTestSuite struct {
	suite.Suite
	grpcServer     *grpcapp.Server
	userGrpcClient userpb.UserServiceClient
	userGrpcConn   *grpc.ClientConn
	userTestUtil   *util.UserTest
	postgresDB     *gorm.DB
	redisDB        *redis.ClusterClient
	conf           *config.Config
	logger         *logrus.Logger
	user           *entity.User
}

func (f *FindUserByEmailTestSuite) SetupSuite() {
	grpcServer, postgresDB, redisDB, conf, logger := util.NewGrpcServer()
	f.grpcServer = grpcServer
	f.postgresDB = postgresDB
	f.redisDB = redisDB
	f.conf = conf
	f.logger = logger

	f.userTestUtil = util.NewUserTest(postgresDB, logger)

	go f.grpcServer.Run()

	time.Sleep(2 * time.Second)

	userGrpcClient, userGrpcConn := util.NewGrpcUserClient(f.conf.ApiGateway.BaseUrl)
	f.userGrpcClient = userGrpcClient
	f.userGrpcConn = userGrpcConn

	f.user = f.userTestUtil.Create()
}

func (f *FindUserByEmailTestSuite) TearDownSuite() {
	f.userTestUtil.Delete()

	f.grpcServer.Stop()
	f.userGrpcConn.Close()

	f.redisDB.Close()
	sqlDB, _ := f.postgresDB.DB()
	sqlDB.Close()
}

func (f *FindUserByEmailTestSuite) Test_Success() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	res, err := f.userGrpcClient.FindByEmail(ctx, &userpb.Email{Email: f.user.Email})

	assert.NoError(f.T(), err)
	assert.NotNil(f.T(), res.Data)

	st, _ := status.FromError(err)
	assert.Equal(f.T(), codes.OK, st.Code())
}

func (f *FindUserByEmailTestSuite) Test_NotFound() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	res, err := f.userGrpcClient.FindByEmail(ctx, &userpb.Email{Email: "usernotfound@gmail.com"})

	assert.NoError(f.T(), err)
	assert.Nil(f.T(), res.Data)
}

func (f *FindUserByEmailTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := f.userGrpcClient.FindByEmail(ctx, &userpb.Email{Email: f.user.Email})

	st, _ := status.FromError(err)
	assert.Equal(f.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_FindUserByEmail(t *testing.T) {
	suite.Run(t, new(FindUserByEmailTestSuite))
}
