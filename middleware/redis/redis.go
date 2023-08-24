package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

var (
	Ctx        = context.Background()
	RdbLike    *redis.Client
	RdbMessage *redis.Client

	// ...
	// 根据需要再添加多个client
)

// InitRedis 初始化Redis
func InitRedis() {
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")

	// 初始化所有client
	// 不同client连接不同的Redis DB
	RdbLike = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1,
	})

	RdbMessage = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2,
	})
}

func WithTimeoutContextBySecond(second time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(Ctx, second*time.Second)
}
