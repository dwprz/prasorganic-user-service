package integration_test

import (
	"context"
	"encoding/base64"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
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
	"testing"
	"time"
)

// *nyalakan nginx dan database nya terlebih dahulu
// go test -v ./test/integration -count=1 -p=1
// go test -run ^TestIntegration_UpsertUser$ -v ./test/integration -count=1

type UpsertUserTestSuite struct {
	suite.Suite
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

func (u *UpsertUserTestSuite) SetupSuite() {
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
}

func (u *UpsertUserTestSuite) TearDownSuite() {
	u.grpcServer.Stop()
	u.userGrpcConn.Close()

	u.redisTestUtil.Flushall()
	u.redisDB.Close()

	sqlDB, _ := u.postgresDB.DB()
	sqlDB.Close()
}

func (u *UpsertUserTestSuite) TearDownTest() {
	u.userTestUtil.Delete()
}

func (u *UpsertUserTestSuite) Test_Success() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &pb.LoginWithGoogleRequest{
		UserId:       "ynA1nZIULkXLrfy0fvz5t",
		Email:        "johndoe123@gmail.com",
		FullName:     "John Doe",
		PhotoProfile: "example-photo-profile",
		RefreshToken: `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOj
					   E3MjUxNzIwMDUsImlkIjoiMV9pUGtNbjk4c19ObXNRZ1Q1T
					   WtlIiwiaXNzIjoicHJhc29yZ2FuaWMtYXV0aC1zZXJ2aWNl
					   In0.cVJL1ivJ5wDECYwBQtA39R_HMkEaG4HiRHxZSJBl0EL
					   5_EcuKq5v7QscveiFYd7CEsRRtnHv3hvosa7pndWgZwfOBY
					   pmAybLh6mfgjADUXxtvBzPMT7NGab2rv5ORiv8y4FvOQ45x
					   eKwNKz0Wr2wxiD4tfyzop3_D9OB-ta3F6E`,
	}

	res, err := u.userGrpcClient.Upsert(ctx, req)
	assert.NoError(u.T(), err)

	assert.Equal(u.T(), req.UserId, res.UserId)
	assert.Equal(u.T(), req.Email, res.Email)
	assert.Equal(u.T(), req.FullName, res.FullName)
	assert.Equal(u.T(), req.PhotoProfile, res.PhotoProfile)
	assert.Equal(u.T(), "USER", res.Role)
	assert.NotEmpty(u.T(), res.CreatedAt)
	assert.Equal(u.T(), req.RefreshToken, res.RefreshToken)

}

func (u *UpsertUserTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.LoginWithGoogleRequest{
		UserId:       "ynA1nZIULkXLrfy0fvz5t",
		Email:        "johndoe123@gmail.com",
		FullName:     "John Doe",
		PhotoProfile: "example-photo-profile",
		RefreshToken: `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOj
					   E3MjUxNzIwMDUsImlkIjoiMV9pUGtNbjk4c19ObXNRZ1Q1T
					   WtlIiwiaXNzIjoicHJhc29yZ2FuaWMtYXV0aC1zZXJ2aWNl
					   In0.cVJL1ivJ5wDECYwBQtA39R_HMkEaG4HiRHxZSJBl0EL
					   5_EcuKq5v7QscveiFYd7CEsRRtnHv3hvosa7pndWgZwfOBY
					   pmAybLh6mfgjADUXxtvBzPMT7NGab2rv5ORiv8y4FvOQ45x
					   eKwNKz0Wr2wxiD4tfyzop3_D9OB-ta3F6E`,
	}

	_, err := u.userGrpcClient.Upsert(ctx, req)
	st, _ := status.FromError(err)
	assert.Equal(u.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_UpsertUser(t *testing.T) {
	suite.Run(t, new(UpsertUserTestSuite))
}
