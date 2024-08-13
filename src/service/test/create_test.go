package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/mock/cache"
	"github.com/dwprz/prasorganic-user-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-user-service/src/mock/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	serviceimpl "github.com/dwprz/prasorganic-user-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// go test -v ./src/service/test/... -count=1 -p=1
// go test -run ^TestService_Create$ -v ./src/service/test -count=1

type CreateTestSuite struct {
	suite.Suite
	userService service.User
	userRepo    *repository.UserMock
	userCache *cache.UserMock
}

func (c *CreateTestSuite) SetupSuite() {
	// mock
	c.userRepo = repository.NewUserMock()
	c.userCache = cache.NewUserMock()
	otpGrpcDelivery := delivery.NewOtpGrpcMock()
	otpGrpcConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(otpGrpcDelivery, otpGrpcConn)
	c.userService = serviceimpl.NewUser(grpcClient, c.userRepo, c.userCache)
}

func (c *CreateTestSuite) Test_Succsess() {
	req := &dto.CreateReq{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	c.userRepo.Mock.On("Create", mock.Anything, req).Return(nil)

	err := c.userService.Create(context.Background(), req)
	assert.NoError(c.T(), err)
}

func (c *CreateTestSuite) Test_InvalidEmail() {
	req := &dto.CreateReq{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "123456",
		FullName: "John Doe",
		Password: "rahasia",
	}
	err := c.userService.Create(context.Background(), req)
	assert.Error(c.T(), err)

	errVldtn, ok := err.(validator.ValidationErrors)
	assert.True(c.T(), ok)

	assert.Equal(c.T(), "Email", errVldtn[0].Field())
}

func (c *CreateTestSuite) Test_AlreadyExists() {
	req := &dto.CreateReq{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "existeduser@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	errorRes := &errors.Response{HttpCode: 409, GrpcCode: codes.AlreadyExists, Message: "user already exists"}
	c.userRepo.Mock.On("Create", mock.Anything, req).Return(errorRes)

	err := c.userService.Create(context.Background(), req)
	assert.Error(c.T(), err)
	assert.Equal(c.T(), errorRes, err)
}

func TestService_Create(t *testing.T) {
	suite.Run(t, new(CreateTestSuite))
}
