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
// go test -run ^TestIntegration_Create$ -v ./test/integration -count=1

type CreateTestSuite struct {
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

func (c *CreateTestSuite) SetupSuite() {
	c.postgresDB = database.NewPostgres()
	c.redisDB = database.NewRedisCluster()

	otpGrpcDelivery := delivery.NewOtpGrpcMock()
	grpcClient := util.InitGrpcClientTest(otpGrpcDelivery)

	userService := util.InitUserServiceTest(grpcClient, c.postgresDB, c.redisDB)
	c.grpcServer = util.InitGrpcServerTest(userService)

	go c.grpcServer.Run()

	time.Sleep(1 * time.Second)

	userGrpcDelivery, userGrpcConn := util.InitUserGrpcDelivery()
	c.userGrpcDelivery = userGrpcDelivery
	c.userGrpcConn = userGrpcConn

	c.userTestUtil = util.NewUserTest(c.postgresDB)
	c.redisTestUtil = util.NewRedisTest(c.redisDB)
}

func (c *CreateTestSuite) TearDownSuite() {
	c.grpcServer.Stop()
	c.userGrpcConn.Close()

	c.redisTestUtil.Flushall()
	c.redisDB.Close()

	sqlDB, _ := c.postgresDB.DB()
	sqlDB.Close()
}

func (c *CreateTestSuite) TearDownTest() {
	c.userTestUtil.Delete()
}

func (c *CreateTestSuite) Test_Success() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &pb.RegisterReq{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "Rahasia",
	}

	_, err := c.userGrpcDelivery.Create(ctx, req)
	assert.NoError(c.T(), err)
}

func (c *CreateTestSuite) Test_WithouthUserId() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &pb.RegisterReq{
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "Rahasia",
	}

	_, err := c.userGrpcDelivery.Create(ctx, req)

	st, _ := status.FromError(err)
	assert.Equal(c.T(), codes.InvalidArgument, st.Code())
}

func (c *CreateTestSuite) Test_AlreadyExists() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &pb.RegisterReq{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "Rahasia",
	}

	c.userGrpcDelivery.Create(ctx, req)
	_, err := c.userGrpcDelivery.Create(ctx, req)

	st, _ := status.FromError(err)
	assert.Equal(c.T(), codes.AlreadyExists, st.Code())
}

func (c *CreateTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.RegisterReq{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "Rahasia",
	}

	_, err := c.userGrpcDelivery.Create(ctx, req)
	st, _ := status.FromError(err)

	assert.Equal(c.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_Create(t *testing.T) {
	suite.Run(t, new(CreateTestSuite))
}
