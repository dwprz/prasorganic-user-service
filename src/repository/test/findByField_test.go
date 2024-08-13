package test

import (
	"context"
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
	"testing"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_FindByField$ -v ./src/repository/test -count=1

type FindByFieldTestSuite struct {
	suite.Suite
	user          *entity.User
	userRepo      repository.User
	postgresDB    *gorm.DB
	userCache     cache.User
	redisDB       *redis.ClusterClient
	userTestUtil  *util.UserTest
	redisTestUtil *util.RedisTest
}

func (f *FindByFieldTestSuite) SetupSuite() {
	f.postgresDB = database.NewPostgres()
	f.redisDB = database.NewRedisCluster()

	f.userCache = cacheimpl.NewUser(f.redisDB)

	f.userRepo = repoimpl.NewUser(f.postgresDB, f.userCache)
	f.userTestUtil = util.NewUserTest(f.postgresDB)
	f.redisTestUtil = util.NewRedisTest(f.redisDB)

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
