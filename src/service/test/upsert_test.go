package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-user-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/mock/cache"
	"github.com/dwprz/prasorganic-user-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-user-service/src/mock/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	serviceimpl "github.com/dwprz/prasorganic-user-service/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -v ./src/service/test/... -count=1 -p=1
// go test -run ^TestService_Upsert$ -v ./src/service/test -count=1

type UpsertTestSuite struct {
	suite.Suite
	userService service.User
	userRepo    *repository.UserMock
	userCache   *cache.UserMock
}

func (u *UpsertTestSuite) SetupSuite() {
	// mock
	u.userRepo = repository.NewUserMock()
	u.userCache = cache.NewUserMock()
	otpGrpcDelivery := delivery.NewOtpGrpcMock()
	otpGrpcConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(otpGrpcDelivery, otpGrpcConn)
	u.userService = serviceimpl.NewUser(grpcClient, u.userRepo, u.userCache)
}

func (u *UpsertTestSuite) Test_Succsess() {
	serviceReq := &dto.UpsertReq{
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

	repoRes := &entity.User{
		UserId:       serviceReq.UserId,
		Email:        serviceReq.Email,
		FullName:     serviceReq.FullName,
		PhotoProfile: serviceReq.PhotoProfile,
		Role:         "USER",
		RefreshToken: serviceReq.RefreshToken,
	}

	u.userRepo.Mock.On("Upsert", mock.Anything, serviceReq).Return(repoRes, nil)

	res, err := u.userService.Upsert(context.Background(), serviceReq)
	assert.NoError(u.T(), err)

	assert.Equal(u.T(), repoRes, res)
}

func TestService_Upsert(t *testing.T) {
	suite.Run(t, new(UpsertTestSuite))
}
