package test

import (
	"context"
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
	"testing"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_FindByField$ -v ./src/repository/test -count=1

type FindByFieldTestSuite struct {
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

func (f *FindByFieldTestSuite) SetupSuite() {
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

func (f *FindByFieldTestSuite) TearDownSuite() {
	f.userTestUtil.Delete()
	sqlDB, _ := f.postgresDB.DB()
	sqlDB.Close()

	f.redisTestUtil.Flushall()
	f.redisDB.Close()
}

func (f *FindByFieldTestSuite) Test_Success() {
	res, err := f.userRepo.FindByFields(context.Background(), &entity.User{Email: f.user.Email})
	assert.NoError(f.T(), err)
	assert.Equal(f.T(), f.user, res)
}

func (f *FindByFieldTestSuite) Test_NotFound() {
	user, err := f.userRepo.FindByFields(context.Background(), &entity.User{Email: "notfounduser@gmail.com"})
	assert.NoError(f.T(), err)
	assert.Nil(f.T(), user)
}

func TestRepository_FindByField(t *testing.T) {
	suite.Run(t, new(FindByFieldTestSuite))
}
