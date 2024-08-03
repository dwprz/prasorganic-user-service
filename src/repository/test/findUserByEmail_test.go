package test

import (
	"context"
	chaceinterface "github.com/dwprz/prasorganic-user-service/src/interface/cache"
	repointerface "github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_FindUserByEmail$ -v ./src/repository/test -count=1

type FindUserByEmailTestSuite struct {
	suite.Suite
	user          *entity.User
	userRepo      repointerface.User
	postgresDB    *gorm.DB
	userCache     chaceinterface.User
	redisDB       *redis.ClusterClient
	logger        *logrus.Logger
	userTestUtil  *util.UserTest
	redisTestUtil *util.RedisTest
}

func (f *FindUserByEmailTestSuite) SetupSuite() {
	f.logger = logger.New()
	conf := config.New("DEVELOPMENT", f.logger)
	f.postgresDB = database.NewPostgres(conf)
	f.redisDB = database.NewRedisCluster(conf)

	f.userCache = cache.NewUser(f.redisDB, f.logger)

	f.userRepo = repository.NewUser(f.postgresDB, f.userCache)
	f.userTestUtil = util.NewUserTest(f.postgresDB, f.logger)
	f.redisTestUtil = util.NewRedisTest(f.redisDB, f.logger)

	f.user = f.userTestUtil.Create()
}

func (f *FindUserByEmailTestSuite) TearDownSuite() {
	f.userTestUtil.Delete()
	sqlDB, _ := f.postgresDB.DB()
	sqlDB.Close()

	f.redisTestUtil.Flushall()
	f.redisDB.Close()
}

func (f *FindUserByEmailTestSuite) Test_Success() {
	user, err := f.userRepo.FindByEmail(context.Background(), f.user.Email)
	assert.NoError(f.T(), err)
	assert.Equal(f.T(), f.user, user)
}

func (f *FindUserByEmailTestSuite) Test_NotFound() {
	user, err := f.userRepo.FindByEmail(context.Background(), "notfounduser@gmail.com")
	assert.NoError(f.T(), err)
	assert.Nil(f.T(), user)
}

func TestRepository_FindUserByEmail(t *testing.T) {
	suite.Run(t, new(FindUserByEmailTestSuite))
}
