package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dwprz/prasorganic-user-service/src/interface/cache"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type UserImpl struct {
	redis  *redis.ClusterClient
	logger *logrus.Logger
}

func NewUser(r *redis.ClusterClient, l *logrus.Logger) cache.User {
	return &UserImpl{
		redis:  r,
		logger: l,
	}
}

func (u *UserImpl) Cache(ctx context.Context, user *entity.User) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		u.logger.WithFields(logrus.Fields{"location": "cache.UserImpl/Cache", "section": "Marshal"}).Error(err.Error())
		return
	}

	key := fmt.Sprintf("user:%s", user.Email)
	const expire = 24 * time.Hour

	if _, err := u.redis.SetEx(ctx, key, string(jsonData), expire).Result(); err != nil {
		u.logger.WithFields(logrus.Fields{"location": "cache.UserImpl/Cache", "section": "Marshal"}).Error(err.Error())
	}
}

func (u *UserImpl) FindByEmail(ctx context.Context, email string) *entity.User {
	res, err := u.redis.Get(ctx, fmt.Sprintf("user:%s", email)).Result()

	if err != nil && err != redis.Nil {
		u.logger.WithFields(logrus.Fields{"location": "cache.UserImpl/FindByEmail", "section": "Get"}).Error(err.Error())
		return nil
	}

	if res == "" {
		return nil
	}

	user := new(entity.User)

	if err := json.Unmarshal([]byte(res), user); err != nil {
		u.logger.WithFields(logrus.Fields{"location": "cache.UserImpl/FindByEmail", "section": "Unmarshal"}).Error(err.Error())
		return nil
	}

	return user
}
