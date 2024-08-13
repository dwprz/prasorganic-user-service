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
// go test -run ^TestIntegration_FindByEmail$  -v ./test/integration -count=1

type FindByEmailTestSuite struct {
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

func (f *FindByEmailTestSuite) SetupSuite() {
	f.postgresDB = database.NewPostgres()
	f.redisDB = database.NewRedisCluster()

	otpGrpcDelivery := delivery.NewOtpGrpcMock()
	grpcClient := util.InitGrpcClientTest(otpGrpcDelivery)

	userService := util.InitUserServiceTest(grpcClient, f.postgresDB, f.redisDB)
	f.grpcServer = util.InitGrpcServerTest(userService)

	go f.grpcServer.Run()

	time.Sleep(1 * time.Second)

	userGrpcDelivery, userGrpcConn := util.InitUserGrpcDelivery()
	f.userGrpcDelivery = userGrpcDelivery
	f.userGrpcConn = userGrpcConn

	f.userTestUtil = util.NewUserTest(f.postgresDB)
	f.redisTestUtil = util.NewRedisTest(f.redisDB)

	f.user = f.userTestUtil.Create()
}

func (f *FindByEmailTestSuite) TearDownSuite() {
	f.redisTestUtil.Flushall()
	f.redisDB.Close()

	f.userTestUtil.Delete()
	sqlDB, _ := f.postgresDB.DB()
	sqlDB.Close()

	f.grpcServer.Stop()
	f.userGrpcConn.Close()
}

func (f *FindByEmailTestSuite) Test_Success() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	res, err := f.userGrpcDelivery.FindByEmail(ctx, &pb.Email{Email: f.user.Email})

	assert.NoError(f.T(), err)
	assert.NotNil(f.T(), res.Data)

	st, _ := status.FromError(err)
	assert.Equal(f.T(), codes.OK, st.Code())
}

func (f *FindByEmailTestSuite) Test_NotFound() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	res, err := f.userGrpcDelivery.FindByEmail(ctx, &pb.Email{Email: "usernotfound@gmail.com"})

	assert.NoError(f.T(), err)
	assert.Nil(f.T(), res.Data)
}

func (f *FindByEmailTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := f.userGrpcDelivery.FindByEmail(ctx, &pb.Email{Email: f.user.Email})

	st, _ := status.FromError(err)
	assert.Equal(f.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_FindByEmail(t *testing.T) {
	suite.Run(t, new(FindByEmailTestSuite))
}
