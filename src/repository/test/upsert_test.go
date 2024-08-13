package test

import (
	"context"
	cacheimpl "github.com/dwprz/prasorganic-user-service/src/cache"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-user-service/src/interface/cache"
	"github.com/dwprz/prasorganic-user-service/src/interface/repository"
	"github.com/dwprz/prasorganic-user-service/src/model/dto"
	repoimpl "github.com/dwprz/prasorganic-user-service/src/repository"
	"github.com/dwprz/prasorganic-user-service/test/util"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_Upsert$ -v ./src/repository/test -count=1

type UpsertTestSuite struct {
	suite.Suite
	userRepo      repository.User
	postgresDB    *gorm.DB
	userCache     cache.User
	redisDB       *redis.ClusterClient
	userTestUtil  *util.UserTest
	redisTestUtil *util.RedisTest
}

func (u *UpsertTestSuite) SetupSuite() {
	u.postgresDB = database.NewPostgres()
	u.redisDB = database.NewRedisCluster()

	u.userCache = cacheimpl.NewUser(u.redisDB)

	u.userRepo = repoimpl.NewUser(u.postgresDB, u.userCache)
	u.userTestUtil = util.NewUserTest(u.postgresDB)
	u.redisTestUtil = util.NewRedisTest(u.redisDB)
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
