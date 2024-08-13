package database

import (
	"time"

	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisCluster() *redis.ClusterClient {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			config.Conf.Redis.AddrNode1,
			config.Conf.Redis.AddrNode2,
			config.Conf.Redis.AddrNode3,
			config.Conf.Redis.AddrNode4,
			config.Conf.Redis.AddrNode5,
			config.Conf.Redis.AddrNode6,
		},
		Password:     config.Conf.Redis.Password,
		DialTimeout:  20 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	})

	return client
}
