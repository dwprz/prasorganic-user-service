package test

import (
	"context"
	"testing"
	chaceinterface "github.com/dwprz/prasorganic-user-service/src/interface/cache"
	repointerface "github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/common/logger"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	"github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_Upsert$ -v ./src/repository/test -count=1

type UpsertTestSuite struct {
	suite.Suite
	userRepo      repointerface.User
	postgresDB    *gorm.DB
	userCache     chaceinterface.User
	redisDB       *redis.ClusterClient
	logger        *logrus.Logger
	userTestUtil  *util.UserTest
	redisTestUtil *util.RedisTest
}

func (u *UpsertTestSuite) SetupSuite() {
	u.logger = logger.New()
	conf := config.New("DEVELOPMENT", u.logger)
	u.postgresDB = database.NewPostgres(conf)
	u.redisDB = database.NewRedisCluster(conf)

	u.userCache = cache.NewUser(u.redisDB, u.logger)

	u.userRepo = repository.NewUser(u.postgresDB, u.userCache)
	u.userTestUtil = util.NewUserTest(u.postgresDB, u.logger)
	u.redisTestUtil = util.NewRedisTest(u.redisDB, u.logger)
}

func (u *UpsertTestSuite) TearDownSuite() {
	u.userTestUtil.Delete()
	sqlDB, _ := u.postgresDB.DB()
	sqlDB.Close()

	u.redisTestUtil.Flushall()
	u.redisDB.Close()
}

func (u *UpsertTestSuite) Test_Success() {
	req := &dto.UpsertReq{
		UserId:       "ynA1nZIULkXLrfy0fvz5t",
		Email:        "johndoe123@gmail.com",
		FullName:     "John Doe",
		PhotoProfile: "example-photo-profile",
		RefreshToken: `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOj
					   E3MjUxNzIwMDUsImlkIjoiMV9pUGtNbjk4c19ObXNRZ1Q1T
					   WtlIiwiaXNzIjoicHJhc29yZ2FuaWMtYXV0aC1zZXJ2aWNl
					   In0.cVJL1ivJ5wDECYwBQtA39R_HMkEaG4HiRHxZSJBl0EL
					   5_EcuKq5v7QscveiFYd7CEsRRtnHv3hvosa7pndWgZwfOBY
					   pmAybLh6mfgjADUXxtvBzPMT7NGab2rv5ORiv8y4FvOQ45x
					   eKwNKz0Wr2wxiD4tfyzop3_D9OB-ta3F6E`,
	}

	res, err := u.userRepo.Upsert(context.Background(), req)
	assert.NoError(u.T(), err)

	assert.Equal(u.T(), req.UserId, res.UserId)
	assert.Equal(u.T(), req.Email, res.Email)
	assert.Equal(u.T(), req.FullName, res.FullName)
	assert.Equal(u.T(), req.PhotoProfile, res.PhotoProfile)
	assert.Equal(u.T(), "USER", res.Role)
	assert.NotEmpty(u.T(), res.CreatedAt)
	assert.Equal(u.T(), req.RefreshToken, res.RefreshToken)
}

func TestRepository_Upsert(t *testing.T) {
	suite.Run(t, new(UpsertTestSuite))
}
