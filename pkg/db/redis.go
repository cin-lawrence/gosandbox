package db

import (
	"context"

	"github.com/cin-lawrence/gosandbox/pkg/config"

	miniredis "github.com/alicebob/miniredis/v2"
	redis "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client = InitRedis()
var RedisContext context.Context = context.TODO()

func InitRedis() *redis.Client {
	if config.Config.Test {
		mr, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		return redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
	}
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisAddress,
		Password: config.Config.RedisPassword,
		DB:       config.Config.RedisDB,
	})
}
