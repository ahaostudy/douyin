package service

import (
	"context"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"sync"
)

func GetFollowerList(id, curID uint) ([]*model.User, bool) {
	FollowerList, err := dao.GetFollowerList(id, curID)
	return FollowerList, err == nil
}

// LoadFollowerList 加载粉丝列表到redis
func LoadFollowerList(ctx context.Context, id uint) error {
	key := redis.GenerateFollowerKey(id)

	// 刷新过期时间
	defer func() {
		go func() { redis.RdbFollower.Expire(ctx, key, config.RedisKeyTTL) }()
	}()

	// 判断redis里是否存在关注列表，不存在时从MySQL读取到Redis
	if n, err := redis.RdbFollower.Exists(ctx, key).Result(); n == 0 || err != nil {
		// 从数据库读取数据
		followerList, err := dao.GetBasicFollowerList(id)
		if err != nil {
			return err
		}

		// 先添加一个默认值，防止缓存穿透
		if err := redis.RdbFollower.SAdd(ctx, key, config.RedisValueOfNULL).Err(); err != nil {
			return err
		}

		// 并发写入redis
		wg := sync.WaitGroup{}
		for _, follower := range followerList {
			wg.Add(1)
			go func(follower model.Follow) {
				defer wg.Done()
				if e := redis.RdbFollower.SAdd(ctx, key, follower.FollowerID).Err(); e != nil {
					err = e
				}
			}(*follower)
		}
		wg.Wait()

		// 如果过程中读取失败，直接将key删除，防止脏写
		if err != nil {
			redis.RdbFollower.Del(ctx, key)
			return err
		}
	}

	return nil
}
