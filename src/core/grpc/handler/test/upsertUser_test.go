package test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-user-service/src/mock/service"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// go test -v ./src/core/grpc/handler/test/... -count=1 -p=1
// go test -run ^TestServer_UpsertUser$ -v ./src/core/grpc/handler/test -count=1

type UpsertUserTestSuite struct {
	suite.Suite
	grpcServer       *server.Grpc
	userGrpcDelivery pb.UserServiceClient
	userGrpcConn     *grpc.ClientConn
	userService      *service.UserMock
}

func (u *UpsertUserTestSuite) SetupSuite() {
	u.userService = service.NewUserMock()
	u.grpcServer = util.InitGrpcServerTest(u.userService)

	go u.grpcServer.Run()

	time.Sleep(1 * time.Second)

	userGrpcDelivery, userGrpcConn := util.InitUserGrpcDelivery()
	u.userGrpcDelivery = userGrpcDelivery
	u.userGrpcConn = userGrpcConn
}

func (u *UpsertUserTestSuite) TearDownSuite() {
	u.grpcServer.Stop()
	u.userGrpcConn.Close()
}

func (u *UpsertUserTestSuite) Test_Success() {
	serverReq := &pb.LoginWithGoogleReq{
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

	serviceReq := new(dto.UpsertReq)
	copier.Copy(serviceReq, serverReq)

	serviceRes := &entity.User{
		UserId:       serverReq.UserId,
		Email:        serverReq.Email,
		FullName:     serverReq.FullName,
		PhotoProfile: serverReq.PhotoProfile,
		Role:         "USER",
		RefreshToken: serverReq.RefreshToken,
	}

	u.userService.Mock.On("Upsert", mock.Anything, serviceReq).Return(serviceRes, nil)

	auth := base64.StdEncoding.EncodeToString([]byte("prasorganic-auth:rahasia"))
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Basic "+auth)

	res, err := u.userGrpcDelivery.Upsert(ctx, serverReq)
	assert.NoError(u.T(), err)

	user := new(pb.User)
	err = copier.Copy(user, res)
	assert.NoError(u.T(), err)

	assert.Equal(u.T(), user, res)
}

func TestServer_UpsertUser(t *testing.T) {
	suite.Run(t, new(UpsertUserTestSuite))
}
