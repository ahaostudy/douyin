package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	Ctx              = context.Background()
	RdbLikeByUserID  *redis.Client
	RdbLikeByVideoID *redis.Client

	// ...
	// 根据需要再添加多个client
)

// InitRedis 初始化Redis
func InitRedis() {
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")

	// 初始化所有client
	// 不同client连接不同的Redis DB
	RdbLikeByUserID = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	RdbLikeByVideoID = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1,
	})

	// ...
}
