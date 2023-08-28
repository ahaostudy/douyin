package service

import (
	"context"
	"encoding/json"
	"fmt"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"sort"
	"sync"
)

// GetCommentList 获取评论列表
func GetCommentList(vid uint, uid uint) ([]*model.Comment, bool) {
	fmt.Println("enter comment list")
	var commentList []*model.Comment
	// var err error
	// 先从redis中加载相关的视频评论信息
	ctx, cancel := redis.WithTimeoutContextBySecond(300)
	defer cancel()
	key := redis.GenerateCommentKey(vid)

	// 1. 加载评论列表到redis
	if LoadCommentList(ctx, vid, uid) != nil {
		return nil, false
	}

	commentHash, err := redis.RdbComment.HGetAll(ctx, key).Result()

	if err != nil {
		return nil, false
	}

	if len(commentHash) != 0 {
		for _, v := range commentHash {
			var commentType model.Comment

			if err := json.Unmarshal([]byte(v), &commentType); err != nil {
				return nil, false
			}
			commentList = append(commentList, &commentType)
		}
		// 处理并发 HSet过程中导致的创建日期乱序情况
		// 排序评论列表
		sort.Slice(commentList, func(i, j int) bool {
			idI := commentList[i].ID
			idJ := commentList[j].ID
			return idI > idJ // 从小到大排序，若要从大到小排序可改为 idI > idJ
		})
	}

	return commentList, true
}

// LoadCommentList 加载评论列表到redis
// 评论是依赖于视频的，所以传入视频id vid
func LoadCommentList(ctx context.Context, vid, uid uint) error {
	key := redis.GenerateCommentKey(vid)

	// 刷新过期时间
	defer func() {
		go func() { redis.RdbComment.Expire(ctx, key, config.RedisKeyTTL) }()
	}()

	// 判断redis里是否存在该视频的评论列表，不存在时从MySQL读取到Redis
	if n, err := redis.RdbComment.Exists(ctx, key).Result(); n == 0 || err != nil {
		// 从数据库读取数据
		commentList, err := dao.GetCommentList(vid, uid)
		if err != nil {
			return err
		}

		// 先添加一个默认值，防止缓存穿透
		if err := redis.RdbComment.HSet(ctx, key, config.RedisValueOfNULL, config.RedisValueOfNULL).Err(); err != nil {
			return err
		}

		// 并发写入redis
		wg := sync.WaitGroup{}
		for _, comment := range commentList {
			wg.Add(1)
			go func(comment model.Comment) {
				defer wg.Done()
				commentJson, err1 := json.Marshal(comment)
				if err != nil {
					err = err1
				}
				if e := redis.RdbComment.HSet(ctx, key, comment.ID, commentJson).Err(); e != nil {
					err = e
				}
			}(*comment)
		}
		wg.Wait()
		// for _, comment := range commentList {
		// 	commentJson, err1 := json.Marshal(comment)
		// 	if err != nil {
		// 		err = err1
		// 	}
		// 	if e := redis.RdbComment.HSet(ctx, key, comment.ID, commentJson).Err(); e != nil {
		// 		err = e
		// 	}
		// }

		// 清除占位防止缓存穿透的NULL key
		existNULL, err := redis.RdbComment.HExists(ctx, key, config.RedisValueOfNULL).Result()
		if err != nil {
			return err
		}

		if existNULL {
			err := redis.RdbComment.HDel(ctx, key, config.RedisValueOfNULL).Err()
			if err != nil {
				return err
			}
		}

		// 如果过程中读取失败，直接将key删除，防止脏写
		if err != nil {
			redis.RdbComment.Del(ctx, key)
			return err
		}
	}

	return nil
}
