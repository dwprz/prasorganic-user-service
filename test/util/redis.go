package util

import (
	"context"
	"strings"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisTest struct {
	redisDB *redis.ClusterClient
	logger  *logrus.Logger
}

func NewRedisTest(r *redis.ClusterClient, l *logrus.Logger) *RedisTest {
	return &RedisTest{
		redisDB: r,
		logger:  l,
	}
}

func (r *RedisTest) Flushall() {
	nodesInfo, err := r.redisDB.ClusterNodes(context.Background()).Result()
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"location": "util.Redis/Flushall",
			"section":  "redisDB.ClusterNodes",
		}).Error(err)
	}

	// proses informasi node
	nodes := strings.Split(nodesInfo, "\n")
	for _, node := range nodes {
		fields := strings.Fields(node)

		if len(fields) > 2 && strings.Contains(fields[2], "master") {
			err := r.redisDB.Do(context.Background(), "FLUSHALL").Err()
			if err != nil {
				r.logger.WithFields(logrus.Fields{
					"location": "util.Redis/Flushall",
					"section":  "redisDB.Do",
				}).Error(err)
			}
		}
	}
}
