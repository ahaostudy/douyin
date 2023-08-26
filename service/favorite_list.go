package service

import (
	"context"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"sync"
)

// GetFavoriteList 获取用户喜欢的视频列表
func GetFavoriteList(id, curID uint) ([]*model.Video, bool) {
	videoList, err := dao.GetVideoListByLike(id, curID)
	if err != nil {
		return nil, false
	}
	return videoList, true
}

// LoadLikeList 读取视频点赞列表到redis
func LoadLikeList(ctx context.Context, id uint) error {
	key := redis.GenerateLikeKey(id)

	// 刷新过期时间
	defer func() {
		go func() { redis.RdbLike.Expire(ctx, key, config.RedisKeyTTL) }()
	}()

	// 判断redis中是否存在视频的喜欢列表，不存在时从MySQL从读取到redis
	if n, err := redis.RdbLike.Exists(ctx, key).Result(); n == 0 || err != nil {
		// 从MySQL中获取数据
		likeList, err := dao.GetLikeListByVideoID(id)
		if err != nil {
			return err
		}

		// 先写入一个空value，这样就算key没有数据，也会保留在redis中，防止缓存穿透
		if err := redis.RdbLike.SAdd(ctx, key, config.RedisValueOfNULL).Err(); err != nil {
			return err
		}

		// 遍历点赞列表并发添加到redis中
		wg := sync.WaitGroup{}
		for _, like := range likeList {
			wg.Add(1)
			go func(like model.Like) {
				defer wg.Done()
				if e := redis.RdbLike.SAdd(ctx, key, like.UserID).Err(); err != nil {
					err = e
				}
			}(*like)
		}
		wg.Wait()

		// 如果过程中读取失败，直接将key删除，防止脏写
		if err != nil {
			redis.RdbLike.Del(ctx, key)
			return err
		}
	}

	return nil
}
