package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	chaceinterface "github.com/dwprz/prasorganic-user-service/src/interface/cache"
	repointerface "github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_UpdateByEmail$ -v ./src/repository/test -count=1

type UpdateByEmailTestSuite struct {
	suite.Suite
	user          entity.User
	userRepo      repointerface.User
	postgresDB    *gorm.DB
	userCache     chaceinterface.User
	redisDB       *redis.ClusterClient
	logger        *logrus.Logger
	userTestUtil  *util.UserTest
	redisTestUtil *util.RedisTest
}

func (u *UpdateByEmailTestSuite) SetupSuite() {
	u.logger = logger.New()
	conf := config.New("DEVELOPMENT", u.logger)
	u.postgresDB = database.NewPostgres(conf)
	u.redisDB = database.NewRedisCluster(conf)

	u.userCache = cache.NewUser(u.redisDB, u.logger)

	u.userRepo = repository.NewUser(u.postgresDB, u.userCache)
	u.userTestUtil = util.NewUserTest(u.postgresDB, u.logger)
	u.redisTestUtil = util.NewRedisTest(u.redisDB, u.logger)

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
