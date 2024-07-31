package database

import (
	"time"

	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisCluster(conf *config.Config) *redis.ClusterClient {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			conf.Redis.AddrNode1,
			conf.Redis.AddrNode2,
			conf.Redis.AddrNode3,
			conf.Redis.AddrNode4,
			conf.Redis.AddrNode5,
			conf.Redis.AddrNode6,
		},
		Password:     conf.Redis.Password,
		DialTimeout:  20 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	})

	return client
}
