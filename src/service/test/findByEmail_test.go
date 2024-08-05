package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-user-service/src/common/helper"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	grpcapp "github.com/dwprz/prasorganic-user-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	svcinterface "github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/mock/cache"
	"github.com/dwprz/prasorganic-user-service/src/mock/client"
	"github.com/dwprz/prasorganic-user-service/src/mock/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -v ./src/service/test/... -count=1 -p=1
// go test -run ^TestService_FindByEmail$ -v ./src/service/test -count=1

type FindByEmailTestSuite struct {
	suite.Suite
	userService svcinterface.User
	userRepo    *repository.UserMock
	userCache   *cache.UserMock
}

func (f *FindByEmailTestSuite) SetupSuite() {
	logger := logger.New()
	conf := config.New("DEVELOPMENT", logger)
	validator := validator.New()

	// mock
	f.userRepo = repository.NewUserMock()
	f.userCache = cache.NewUserMock()
	otpGrpcClient := client.NewOtpGrpcMock()
	otpGrpcConn := new(grpc.ClientConn)

	grpcClient := grpcapp.NewClient(otpGrpcClient, otpGrpcConn, logger)
	helper := helper.New(conf, logger)
	f.userService = service.NewUser(grpcClient, validator, f.userRepo, f.userCache, helper)
}

func (f *FindByEmailTestSuite) Test_Succsess() {
	req := &entity.User{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
	}

	f.userCache.Mock.On("FindByEmail", mock.Anything, req.Email).Return(req)

	f.userRepo.Mock.On("FindByFields", mock.Anything, &entity.User{
		Email: req.Email,
	}).Return(req, nil)

	res, err := f.userService.FindByEmail(context.Background(), req.Email)
	assert.NoError(f.T(), err)
	assert.Equal(f.T(), req, res)
}

func (f *FindByEmailTestSuite) Test_NotFound() {
	email := "notfounduser@gmail.com"

	f.userCache.Mock.On("FindByEmail", mock.Anything, email).Return(nil)
	f.userRepo.Mock.On("FindByFields", mock.Anything, &entity.User{
		Email: email,
	}).Return(nil, nil)

	res, err := f.userService.FindByEmail(context.Background(), email)
	assert.NoError(f.T(), err)
	assert.Nil(f.T(), res)
}

func TestService_FindByEmail(t *testing.T) {
	suite.Run(t, new(FindByEmailTestSuite))
}
