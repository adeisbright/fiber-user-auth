package loaders

import (
	"github.com/go-redis/redis"
)

func ConnectToRedis() *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return redisClient
}
