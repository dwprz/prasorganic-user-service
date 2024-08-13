package integration_test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/redis/go-redis/v9"
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
// go test -run ^TestIntegration_SetNullRefreshToken$ -v ./test/integration -count=1

type SetNullRefreshTokenTestSuite struct {
	suite.Suite
	user             *entity.User
	grpcServer       *server.Grpc
	userGrpcDelivery pb.UserServiceClient
	userGrpcConn     *grpc.ClientConn
	userTestUtil     *util.UserTest
	postgresDB       *gorm.DB
	redisDB          *redis.ClusterClient
	redisTestUtil    *util.RedisTest
}

func (s *SetNullRefreshTokenTestSuite) SetupSuite() {
	s.postgresDB = database.NewPostgres()
	s.redisDB = database.NewRedisCluster()

	otpGrpcDelivery := delivery.NewOtpGrpcMock()
	grpcClient := util.InitGrpcClientTest(otpGrpcDelivery)

	userService := util.InitUserServiceTest(grpcClient, s.postgresDB, s.redisDB)
	s.grpcServer = util.InitGrpcServerTest(userService)

	go s.grpcServer.Run()

	time.Sleep(1 * time.Second)

	userGrpcDelivery, userGrpcConn := util.InitUserGrpcDelivery()
	s.userGrpcDelivery = userGrpcDelivery
	s.userGrpcConn = userGrpcConn

	s.userTestUtil = util.NewUserTest(s.postgresDB)
	s.redisTestUtil = util.NewRedisTest(s.redisDB)

	s.user = s.userTestUtil.Create()
}

func (s *SetNullRefreshTokenTestSuite) TearDownSuite() {
	s.userTestUtil.Delete()
	sqlDB, _ := s.postgresDB.DB()
	sqlDB.Close()

	s.redisTestUtil.Flushall()
	s.redisDB.Close()

	s.grpcServer.Stop()
	s.userGrpcConn.Close()
}

func (s *SetNullRefreshTokenTestSuite) Test_SetNull() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &pb.RefreshToken{Token: s.user.RefreshToken}

	_, err := s.userGrpcDelivery.SetNullRefreshToken(ctx, req)
	assert.NoError(s.T(), err)
}

func (s *SetNullRefreshTokenTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.RefreshToken{
		Token: s.user.RefreshToken,
	}

	_, err := s.userGrpcDelivery.SetNullRefreshToken(ctx, req)
	st, _ := status.FromError(err)
	assert.Equal(s.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_SetNullRefreshToken(t *testing.T) {
	suite.Run(t, new(SetNullRefreshTokenTestSuite))
}
