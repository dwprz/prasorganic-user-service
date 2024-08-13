package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-user-service/src/interface/cache"
	"github.com/dwprz/prasorganic-user-service/src/interface/repository"
	cacheimpl "github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	repoimpl"github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_Create$ -v ./src/repository/test -count=1

type CreateTestSuite struct {
	suite.Suite
	userRepo      repository.User
	postgresDB    *gorm.DB
	userCache     cache.User
	redisDB       *redis.ClusterClient
	userTestUtil  *util.UserTest
	redisTestUtil *util.RedisTest
}

func (c *CreateTestSuite) SetupSuite() {
	c.postgresDB = database.NewPostgres()
	c.redisDB = database.NewRedisCluster()

	c.userCache = cacheimpl.NewUser(c.redisDB)

	c.userRepo = repoimpl.NewUser(c.postgresDB, c.userCache)
	c.userTestUtil = util.NewUserTest(c.postgresDB)
	c.redisTestUtil = util.NewRedisTest(c.redisDB)
}

func (c *CreateTestSuite) TearDownTest() {
	c.userTestUtil.Delete()

}

func (c *CreateTestSuite) TearDownSuite() {
	sqlDB, _ := c.postgresDB.DB()
	sqlDB.Close()

	c.redisTestUtil.Flushall()
	c.redisDB.Close()
}

func (c *CreateTestSuite) Test_Success() {
	user := &dto.CreateReq{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
		Email:    "johndoe@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	err := c.userRepo.Create(context.Background(), user)
	assert.NoError(c.T(), err)
}

func (c *CreateTestSuite) Test_AlreadyExists() {
	user := &dto.CreateReq{
		UserId:   "ynA1nZIULkXLrfy0fvz5t",
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

func TestRepository_Create(t *testing.T) {
	suite.Run(t, new(CreateTestSuite))
}
