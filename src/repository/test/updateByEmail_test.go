package test

import (
	"context"
	"testing"

	cacheimpl "github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/interface/cache"
	"github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	repoimpl "github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_UpdateByEmail$ -v ./src/repository/test -count=1

type UpdateByEmailTestSuite struct {
	suite.Suite
	user          entity.User
	userRepo      repository.User
	postgresDB    *gorm.DB
	userCache     cache.User
	redisDB       *redis.ClusterClient
	userTestUtil  *util.UserTest
	redisTestUtil *util.RedisTest
}

func (u *UpdateByEmailTestSuite) SetupSuite() {
	u.postgresDB = database.NewPostgres()
	u.redisDB = database.NewRedisCluster()

	u.userCache = cacheimpl.NewUser(u.redisDB)

	u.userRepo = repoimpl.NewUser(u.postgresDB, u.userCache)
	u.userTestUtil = util.NewUserTest(u.postgresDB)
	u.redisTestUtil = util.NewRedisTest(u.redisDB)

	u.user = *u.userTestUtil.Create()
}

func (u *UpdateByEmailTestSuite) TearDownSuite() {
	u.userTestUtil.Delete()
	sqlDB, _ := u.postgresDB.DB()
	sqlDB.Close()

	u.redisTestUtil.Flushall()
	u.redisDB.Close()
}

func (u *UpdateByEmailTestSuite) Test_Success() {
	req := &entity.User{
		Email:    u.user.Email,
		FullName: "new full name",
	}

	res, err := u.userRepo.UpdateByEmail(context.Background(), req)
	assert.NoError(u.T(), err)

	assert.Equal(u.T(), u.user.UserId, res.UserId)
	assert.Equal(u.T(), u.user.Email, res.Email)
	assert.Equal(u.T(), req.FullName, res.FullName)
	assert.Equal(u.T(), u.user.PhotoProfile, res.PhotoProfile)
	assert.Equal(u.T(), "USER", res.Role)
	assert.NotEmpty(u.T(), res.CreatedAt)
	assert.NotEmpty(u.T(), res.UpdatedAt)
	assert.Equal(u.T(), u.user.RefreshToken, res.RefreshToken)
}

func TestRepository_UpdateByEmail(t *testing.T) {
	suite.Run(t, new(UpdateByEmailTestSuite))
}
