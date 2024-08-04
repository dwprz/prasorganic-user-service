package integration_test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
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

// *nyalakan nginx dan database nya terlebih dahulu
// go test -v ./test/integration -count=1 -p=1
// go test -run ^TestIntegration_AddRefreshToken$ -v ./test/integration -count=1

type AddRefreshTokenTestSuite struct {
	suite.Suite
	user           *entity.User
	grpcServer     *grpcapp.Server
	userGrpcClient pb.UserServiceClient
	userGrpcConn   *grpc.ClientConn
	userTestUtil   *util.UserTest
	postgresDB     *gorm.DB
	redisDB        *redis.ClusterClient
	redisTestUtil  *util.RedisTest
	conf           *config.Config
	logger         *logrus.Logger
}

func (u *AddRefreshTokenTestSuite) SetupSuite() {
	grpcServer, postgresDb, redisDB, conf, logger := util.NewGrpcServer()
	u.grpcServer = grpcServer
	u.postgresDB = postgresDb
	u.redisDB = redisDB
	u.conf = conf
	u.logger = logger

	u.userTestUtil = util.NewUserTest(postgresDb, logger)
	u.redisTestUtil = util.NewRedisTest(u.redisDB, u.logger)

	go u.grpcServer.Run()

	time.Sleep(2 * time.Second)

	userGrpcClient, userGrpcConn := util.NewGrpcUserClient(u.conf.ApiGateway.BaseUrl)
	u.userGrpcClient = userGrpcClient
	u.userGrpcConn = userGrpcConn

	u.user = u.userTestUtil.Create()
}

func (u *AddRefreshTokenTestSuite) TearDownSuite() {
	u.userTestUtil.Delete()
	sqlDB, _ := u.postgresDB.DB()
	sqlDB.Close()

	u.redisTestUtil.Flushall()
	u.redisDB.Close()

	u.grpcServer.Stop()
	u.userGrpcConn.Close()
}

func (u *AddRefreshTokenTestSuite) Test_Success() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &pb.AddRefreshToken{
		Email: u.user.Email,
		Token: u.user.RefreshToken,
	}

	_, err := u.userGrpcClient.AddRefreshToken(ctx, req)
	assert.NoError(u.T(), err)
}

func (u *AddRefreshTokenTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.AddRefreshToken{
		Email: u.user.Email,
		Token: u.user.RefreshToken,
	}

	_, err := u.userGrpcClient.AddRefreshToken(ctx, req)
	st, _ := status.FromError(err)
	assert.Equal(u.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_AddRefreshToken(t *testing.T) {
	suite.Run(t, new(AddRefreshTokenTestSuite))
}
