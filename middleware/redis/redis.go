package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"main/config"
	"time"
)

var (
	Ctx         = context.Background()
	RdbLike     *redis.Client
	RdbOpus     *redis.Client
	RdbAuthor   *redis.Client
	RdbFollow   *redis.Client
	RdbFollower *redis.Client
	RdbUser     *redis.Client
	RdbMessage  *redis.Client
	RdbComment  *redis.Client
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
	RdbOpus = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2,
	})
	RdbAuthor = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       3,
	})
	RdbFollow = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       4,
	})
	RdbFollower = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       5,
	})
	RdbUser = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       6,
	})
	RdbMessage = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       7,
	})
	RdbComment = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       8,
	})
}

func WithTimeoutContextBySecond(second time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(Ctx, second*time.Second)
}

// Lock 获取乐观锁
func Lock(rdb *redis.Client, key string) (string, error) {
	ctx, cancel := WithTimeoutContextBySecond(2)
	defer cancel()

	// 生成锁的key和ID
	lockKey := fmt.Sprintf("%s:%s", config.RedisKeyLock, key)
	lockID := uuid.New().String()

	// Redis 主动轮询取锁
	// 每隔50ms尝试一次取锁，最多取锁5次，取锁成功后返回锁ID
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for i := 0; i < 5; i++ {
		result, err := rdb.SetNX(ctx, lockKey, lockID, time.Second).Result()
		if err == nil && result {
			return lockID, nil
		}
		<-ticker.C
	}

	// 多次取锁失败
	return lockID, errors.New("lock is acquired by others")
}

// Unlock 解锁
func Unlock(rdb *redis.Client, key string, lockID string) bool {
	ctx, cancel := WithTimeoutContextBySecond(2)
	defer cancel()

	lockKey := fmt.Sprintf("%s:%s", config.RedisKeyLock, key)

	// 获取当前锁的ID
	// 检测锁ID是否匹配，不匹配则解锁失败
	lockVal, err := rdb.Get(ctx, lockKey).Result()
	if err != nil || lockVal != lockID {
		return false
	}

	// 解锁
	return rdb.Del(ctx, lockKey).Err() == nil
}
