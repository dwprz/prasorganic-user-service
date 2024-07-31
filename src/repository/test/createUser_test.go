package test

import (
	"context"
	"testing"

	repointerface "github.com/dwprz/prasorganic-user-service/interface/repository"
	"github.com/dwprz/prasorganic-user-service/mock/cache"
	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_CreateUser$ -v ./src/repository/test -count=1

type CreateUserTestSuite struct {
	suite.Suite
	userRepo     repointerface.User
	postgresDB   *gorm.DB
	userCache    *cache.UserMock
	logger       *logrus.Logger
	userTestUtil *util.UserTest
}

func (c *CreateUserTestSuite) SetupSuite() {
	c.logger = logger.New()
	conf := config.New("DEVELOPMENT", c.logger)
	c.postgresDB = database.NewPostgres(conf)

	// mock
	c.userCache = cache.NewUserMock()

	c.userRepo = repository.NewUser(c.postgresDB, c.userCache)
	c.userTestUtil = util.NewUserTest(c.postgresDB, c.logger)
}

func (c *CreateUserTestSuite) TearDownTest() {
	c.userTestUtil.Delete()

}

func (c *CreateUserTestSuite) TearDownSuite() {
	sqlDB, _ := c.postgresDB.DB()
	sqlDB.Close()
}

func (c *CreateUserTestSuite) Test_Success() {
	user := &entity.User{
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	err := c.userRepo.Create(context.Background(), user)
	assert.NoError(c.T(), err)
}

func (c *CreateUserTestSuite) Test_AlreadyExists() {
	user := &entity.User{
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	c.userRepo.Create(context.Background(), user)

	err := c.userRepo.Create(context.Background(), user)
	assert.Error(c.T(), err)

	errorRes := &errors.Response{HttpCode: 409, GrpcCode: codes.AlreadyExists, Message: "user already exists"}
	assert.Equal(c.T(), errorRes, err)
}

func TestRepository_CreateUser(t *testing.T) {
	suite.Run(t, new(CreateUserTestSuite))
}
