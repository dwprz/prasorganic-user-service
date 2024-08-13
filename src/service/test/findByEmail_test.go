package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-user-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/mock/cache"
	"github.com/dwprz/prasorganic-user-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-user-service/src/mock/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	serviceimpl "github.com/dwprz/prasorganic-user-service/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -v ./src/service/test/... -count=1 -p=1
// go test -run ^TestService_FindByEmail$ -v ./src/service/test -count=1

type FindByEmailTestSuite struct {
	suite.Suite
	userService service.User
	userRepo    *repository.UserMock
	userCache   *cache.UserMock
}

func (f *FindByEmailTestSuite) SetupSuite() {
	// mock
	f.userRepo = repository.NewUserMock()
	f.userCache = cache.NewUserMock()
	otpGrpcDelivery := delivery.NewOtpGrpcMock()
	otpGrpcConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(otpGrpcDelivery, otpGrpcConn)
	f.userService = serviceimpl.NewUser(grpcClient, f.userRepo, f.userCache)
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
