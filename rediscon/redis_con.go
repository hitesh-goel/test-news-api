package rediscon

import (
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

//ConnectToRedis returns a connection to redis
func ConnectToRedis() (*redis.Client, error) {
	client = redis.NewClient(&redis.Options{
		Addr:         "redis:6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

//GetRedisService get already connected redis
func GetRedisService() *redis.Client {
	return client
}
