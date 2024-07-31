package test

import (
	"context"
	"testing"
	svcinterface "github.com/dwprz/prasorganic-user-service/interface/service"
	"github.com/dwprz/prasorganic-user-service/mock/cache"
	"github.com/dwprz/prasorganic-user-service/mock/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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
	user := &entity.User{
		UserID:   1,
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
	}

	f.userCache.Mock.On("FindByEmail", mock.Anything, user.Email).Return(user)
	f.userRepo.Mock.On("FindByEmail", mock.Anything, user.Email).Return(user, nil)

	res, err := f.userService.FindByEmail(context.Background(), user.Email)
	assert.NoError(f.T(), err)
	assert.Equal(f.T(), user, res)
}

func (f *FindUserByEmailTestSuite) Test_NotFound() {
	email:= "notfounduser@gmail.com"

	f.userCache.Mock.On("FindByEmail", mock.Anything, email).Return(nil)
	f.userRepo.Mock.On("FindByEmail", mock.Anything, email).Return(nil, nil)

	res, err := f.userService.FindByEmail(context.Background(), email)
	assert.NoError(f.T(), err)
	assert.Nil(f.T(), res)
}

func TestService_FindUserByEmail(t *testing.T) {
	suite.Run(t, new(FindUserByEmailTestSuite))
}
