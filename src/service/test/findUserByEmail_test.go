package test

import (
	"context"
	svcinterface "github.com/dwprz/prasorganic-user-service/src/interface/service"
	"github.com/dwprz/prasorganic-user-service/src/mock/cache"
	"github.com/dwprz/prasorganic-user-service/src/mock/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// go test -v ./src/service/test/... -count=1 -p=1
// go test -run ^TestService_FindUserByEmail$ -v ./src/service/test -count=1

type FindUserByEmailTestSuite struct {
	suite.Suite
	userService svcinterface.User
	userRepo    *repository.UserMock
	userCache   *cache.UserMock
}

func (f *FindUserByEmailTestSuite) SetupSuite() {
	validator := validator.New()

	// mock
	f.userRepo = repository.NewUserMock()
	// mock
	f.userCache = cache.NewUserMock()

	f.userService = service.NewUser(validator, f.userRepo, f.userCache)
}

func (f *FindUserByEmailTestSuite) Test_Succsess() {
	req := &entity.User{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
	}

	f.userCache.Mock.On("FindByEmail", mock.Anything, req.Email).Return(req)
	f.userRepo.Mock.On("FindByEmail", mock.Anything, req.Email).Return(req, nil)

	res, err := f.userService.FindByEmail(context.Background(), req.Email)
	assert.NoError(f.T(), err)
	assert.Equal(f.T(), req, res)
}

func (f *FindUserByEmailTestSuite) Test_NotFound() {
	email := "notfounduser@gmail.com"

	f.userCache.Mock.On("FindByEmail", mock.Anything, email).Return(nil)
	f.userRepo.Mock.On("FindByEmail", mock.Anything, email).Return(nil, nil)

	res, err := f.userService.FindByEmail(context.Background(), email)
	assert.NoError(f.T(), err)
	assert.Nil(f.T(), res)
}

func TestService_FindUserByEmail(t *testing.T) {
	suite.Run(t, new(FindUserByEmailTestSuite))
}
