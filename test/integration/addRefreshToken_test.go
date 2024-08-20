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
// go test -run ^TestIntegration_AddRefreshToken$ -v ./test/integration -count=1

type AddRefreshTokenTestSuite struct {
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

func (a *AddRefreshTokenTestSuite) SetupSuite() {
	a.postgresDB = database.NewPostgres()
	a.redisDB = database.NewRedisCluster()

	otpGrpcDelivery := delivery.NewOtpGrpcMock()
	grpcClient := util.InitGrpcClientTest(otpGrpcDelivery)

	userService := util.InitUserServiceTest(grpcClient, a.postgresDB, a.redisDB)
	a.grpcServer = util.InitGrpcServerTest(userService)

	go a.grpcServer.Run()

	time.Sleep(1 * time.Second)

	userGrpcDelivery, userGrpcConn := util.InitUserGrpcDelivery()
	a.userGrpcDelivery = userGrpcDelivery
	a.userGrpcConn = userGrpcConn

	a.userTestUtil = util.NewUserTest(a.postgresDB)
	a.redisTestUtil = util.NewRedisTest(a.redisDB)

	a.user = a.userTestUtil.Create()
}

func (a *AddRefreshTokenTestSuite) TearDownSuite() {
	a.userTestUtil.Delete()
	sqlDB, _ := a.postgresDB.DB()
	sqlDB.Close()

	a.redisTestUtil.Flushall()
	a.redisDB.Close()

	a.grpcServer.Stop()
	a.userGrpcConn.Close()
}

func (a *AddRefreshTokenTestSuite) Test_Success() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &pb.AddRefreshTokenReq{
		Email: a.user.Email,
		Token: a.user.RefreshToken,
	}

	_, err := a.userGrpcDelivery.AddRefreshToken(ctx, req)
	assert.NoError(a.T(), err)
}

func (a *AddRefreshTokenTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.AddRefreshTokenReq{
		Email: a.user.Email,
		Token: a.user.RefreshToken,
	}

	_, err := a.userGrpcDelivery.AddRefreshToken(ctx, req)
	st, _ := status.FromError(err)
	assert.Equal(a.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_AddRefreshToken(t *testing.T) {
	suite.Run(t, new(AddRefreshTokenTestSuite))
}
