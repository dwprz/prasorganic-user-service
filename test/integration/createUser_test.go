package integration_test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	userpb "github.com/dwprz/prasorganic-proto/protogen/user"
	grpcapp "github.com/dwprz/prasorganic-user-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
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

// go test -run ^TestIntegration_CreateUser$ -v ./test/integration -count=1

type CreateUserTestSuite struct {
	suite.Suite
	grpcServer     *grpcapp.Server
	userGrpcClient userpb.UserServiceClient
	userGrpcConn   *grpc.ClientConn
	userTestUtil   *util.UserTest
	postgresDB     *gorm.DB
	redisDB        *redis.ClusterClient
	conf           *config.Config
	logger         *logrus.Logger
}

func (c *CreateUserTestSuite) SetupSuite() {
	grpcServer, postgresDb, redisDB, conf, logger := util.NewGrpcServer()
	c.grpcServer = grpcServer
	c.postgresDB = postgresDb
	c.redisDB = redisDB
	c.conf = conf
	c.logger = logger

	c.userTestUtil = util.NewUserTest(postgresDb, logger)

	go c.grpcServer.Run()
	
	time.Sleep(2 * time.Second)

	userGrpcClient, userGrpcConn := util.NewGrpcUserClient(c.conf.ApiGateway.BaseUrl)
	c.userGrpcClient = userGrpcClient
	c.userGrpcConn = userGrpcConn
}

func (c *CreateUserTestSuite) TearDownSuite() {
	c.grpcServer.Stop()
	c.userGrpcConn.Close()

	c.redisDB.Close()
	sqlDB, _ := c.postgresDB.DB()
	sqlDB.Close()
}

func (c *CreateUserTestSuite) TearDownTest() {
	c.userTestUtil.Delete()
}

func (c *CreateUserTestSuite) Test_Success() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &userpb.RegisterRequest{
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "Rahasia",
	}

	_, err := c.userGrpcClient.Create(ctx, req)
	assert.NoError(c.T(), err)
}

func (c *CreateUserTestSuite) Test_AlreadyExists() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &userpb.RegisterRequest{
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "Rahasia",
	}

	c.userGrpcClient.Create(ctx, req)
	_, err := c.userGrpcClient.Create(ctx, req)

	st, _ := status.FromError(err)

	assert.Equal(c.T(), codes.AlreadyExists, st.Code())
}

func (c *CreateUserTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &userpb.RegisterRequest{
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "Rahasia",
	}

	_, err := c.userGrpcClient.Create(ctx, req)
	st, _ := status.FromError(err)

	assert.Equal(c.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_CreateUser(t *testing.T) {
	suite.Run(t, new(CreateUserTestSuite))
}
