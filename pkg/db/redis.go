package db

import (
	"context"

	"github.com/cin-lawrence/gosandbox/pkg/config"

	redis "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client = InitRedis()
var RedisContext context.Context = context.TODO()

func InitRedis() *redis.Client {
        return redis.NewClient(&redis.Options{
                Addr: config.Config.RedisAddress,
                Password: config.Config.RedisPassword,
                DB: config.Config.RedisDB,
        })
}
