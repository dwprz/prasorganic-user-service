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
// go test -run ^TestIntegration_Upsert$ -v ./test/integration -count=1

type UpsertTestSuite struct {
	suite.Suite
	grpcServer       *server.Grpc
	userGrpcDelivery pb.UserServiceClient
	userGrpcConn     *grpc.ClientConn
	userTestUtil     *util.UserTest
	postgresDB       *gorm.DB
	redisDB          *redis.ClusterClient
	redisTestUtil    *util.RedisTest
}

func (u *UpsertTestSuite) SetupSuite() {
	u.postgresDB = database.NewPostgres()
	u.redisDB = database.NewRedisCluster()

	otpGrpcDelivery := delivery.NewOtpGrpcMock()
	grpcClient := util.InitGrpcClientTest(otpGrpcDelivery)

	userService := util.InitUserServiceTest(grpcClient, u.postgresDB, u.redisDB)
	u.grpcServer = util.InitGrpcServerTest(userService)

	go u.grpcServer.Run()

	time.Sleep(1 * time.Second)

	userGrpcDelivery, userGrpcConn := util.InitUserGrpcDelivery()
	u.userGrpcDelivery = userGrpcDelivery
	u.userGrpcConn = userGrpcConn

	u.userTestUtil = util.NewUserTest(u.postgresDB)
	u.redisTestUtil = util.NewRedisTest(u.redisDB)
}

func (u *UpsertTestSuite) TearDownSuite() {
	u.grpcServer.Stop()
	u.userGrpcConn.Close()

	u.redisTestUtil.Flushall()
	u.redisDB.Close()

	sqlDB, _ := u.postgresDB.DB()
	sqlDB.Close()
}

func (u *UpsertTestSuite) TearDownTest() {
	u.userTestUtil.Delete()
}

func (u *UpsertTestSuite) Test_Success() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	req := &pb.LoginWithGoogleReq{
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

	res, err := u.userGrpcDelivery.Upsert(ctx, req)
	assert.NoError(u.T(), err)

	assert.Equal(u.T(), req.UserId, res.UserId)
	assert.Equal(u.T(), req.Email, res.Email)
	assert.Equal(u.T(), req.FullName, res.FullName)
	assert.Equal(u.T(), req.PhotoProfile, res.PhotoProfile)
	assert.Equal(u.T(), "USER", res.Role)
	assert.NotEmpty(u.T(), res.CreatedAt)
	assert.Equal(u.T(), req.RefreshToken, res.RefreshToken)

}

func (u *UpsertTestSuite) Test_Unauthenticated() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.LoginWithGoogleReq{
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

	_, err := u.userGrpcDelivery.Upsert(ctx, req)
	st, _ := status.FromError(err)
	assert.Equal(u.T(), codes.Unauthenticated, st.Code())
}

func TestIntegration_Upsert(t *testing.T) {
	suite.Run(t, new(UpsertTestSuite))
}
