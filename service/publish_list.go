package service

import (
	"context"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"strconv"
)

// GetWorkList 获取用户的作品列表
func GetWorkList(authorID uint, curID uint) ([]*model.Video, bool) {
	videoList, err := dao.GetVideoListByAuthorID(authorID, curID)
	if err != nil {
		return nil, false
	}

	return videoList, true
}

// GetAuthorID 获取视频作者ID
func GetAuthorID(ctx context.Context, id uint) (uint, error) {
	key := redis.GenerateAuthorKey(id)

	// 刷新过期时间
	defer func() {
		go func() { redis.RdbAuthor.Expire(ctx, key, config.RedisKeyTTL) }()
	}()

	// 从redis中获取视频的作者ID，如果不存在时从MySQL中读取
	val, err := redis.RdbAuthor.Get(ctx, key).Result()
	if err != nil {
		// 从数据库获取数据
		video, err := dao.GetBasicVideo(id)
		if err != nil {
			return 0, err
		}

		// 写入redis
		redis.RdbAuthor.Set(ctx, key, video.AuthorID, config.RedisKeyTTL)

		// 返回结果
		return video.AuthorID, nil
	}

	// 解析作者ID并返回
	aid, err := strconv.ParseUint(val, 10, 32)
	return uint(aid), err
}
