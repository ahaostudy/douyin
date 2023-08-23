package service

import (
	"fmt"
	"main/dao"
	"main/middleware/rabbitmq"
	"main/middleware/redis"
	"main/model"
	"sync"
	"time"
)

// FavoriteAction 赞业务
// uid, vid 分别对应用户ID和视频ID
// t 表示操作类型
func FavoriteAction(uid, vid uint, t int) bool {
	switch t {
	case 1:
		return like(uid, vid)
	case 2:
		return unLike(uid, vid)
	default:
		return false
	}
}

// 点赞业务
// TODO:
// 1. 第一步需要避免缓存穿透（使用bitmap或者布隆过滤器）
// 2. 从MySQL读取到Redis第二个set中失败时，怎么保证两个set数据一致
// 3. 使用MQ异步更新MySQL，怎么保证Redis和MySQL数据一致
func like(uid, vid uint) bool {
	uidKey, vidKey := generateLikeKeyByUser(uid), generateLikeKeyByVideo(vid)
	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()

	// 1. 判断redis中是否有用户以及视频的点赞列表
	//    不存在时先从MySQL读取列表到redis
	if n, err := redis.RdbLike.Exists(ctx, uidKey).Result(); n == 0 || err != nil {
		// 若读取失败，将该key从redis中删除，防止读取到一半，有数据但不全的情况
		if loadLikeListByUser(uid) != nil {
			redis.RdbLike.Del(ctx, uidKey)
			return false
		}
	}
	if n, err := redis.RdbLike.Exists(ctx, vidKey).Result(); n == 0 || err != nil {
		// 若读取失败，将该key从redis中删除，防止读取到一半，有数据但不全的情况
		if loadLikeListByVideo(vid) != nil {
			redis.RdbLike.Del(ctx, vidKey)
			return false
		}
	}

	// 2. 更新redis数据，进行点赞操作
	//    如果已经有点赞记录或更新失败则返回false
	//    如果更新第二个set时更新失败，需要将另一个set中对应的key也删除，确保两个set一致
	//    更新成功后顺便更新key过期时间，保证热点key不过期
	if n, err := redis.RdbLike.SAdd(ctx, uidKey, vid).Result(); n == 0 || err != nil {
		return false
	}
	if n, err := redis.RdbLike.SAdd(ctx, vidKey, uid).Result(); n == 0 {
		return false
	} else if err != nil {
		redis.RdbLike.SRem(ctx, uidKey, vid)
		return false
	}
	// 更新过期时间
	redis.RdbLike.Expire(ctx, uidKey, 24*time.Hour)
	redis.RdbLike.Expire(ctx, vidKey, 24*time.Hour)

	// 3. 写入redis成功，使用MQ异步更新MySQL
	if err := rabbitmq.RMQLikeAdd.Publish([]byte(fmt.Sprintf("%d %d", uid, vid))); err != nil {
		// 如果消息发送到MQ失败，则撤销redis的数据
		redis.RdbLike.SRem(ctx, uidKey, vid)
		redis.RdbLike.SRem(ctx, vidKey, uid)
		return false
	}

	return true
}

// 取消点赞
// 取消点赞逻辑与点赞逻辑类似，不重复写注释
// TODO: 与like和todolist一样
func unLike(uid, vid uint) bool {
	uidKey, vidKey := generateLikeKeyByUser(uid), generateLikeKeyByVideo(vid)
	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()

	if n, err := redis.RdbLike.Exists(ctx, uidKey).Result(); n == 0 || err != nil {
		if loadLikeListByUser(uid) != nil {
			redis.RdbLike.Del(ctx, uidKey)
			return false
		}
	}
	if n, err := redis.RdbLike.Exists(ctx, vidKey).Result(); n == 0 || err != nil {
		if loadLikeListByVideo(vid) != nil {
			redis.RdbLike.Del(ctx, vidKey)
			return false
		}
	}

	if n, err := redis.RdbLike.SRem(ctx, uidKey, vid).Result(); n == 0 || err != nil {
		return false
	}
	if n, err := redis.RdbLike.SRem(ctx, vidKey, uid).Result(); n == 0 {
		return false
	} else if err != nil {
		redis.RdbLike.SAdd(ctx, uidKey, vid)
		return false
	}
	redis.RdbLike.Expire(ctx, uidKey, 24*time.Hour)
	redis.RdbLike.Expire(ctx, vidKey, 24*time.Hour)

	if err := rabbitmq.RMQLikeDel.Publish([]byte(fmt.Sprintf("%d %d", uid, vid))); err != nil {
		redis.RdbLike.SAdd(ctx, uidKey, vid)
		redis.RdbLike.SAdd(ctx, vidKey, uid)
		return false
	}

	return true
}

// 通过user_id生成key
func generateLikeKeyByUser(id uint) string {
	return fmt.Sprintf("like_user:%d", id)
}

// 通过video_id生成key
func generateLikeKeyByVideo(id uint) string {
	return fmt.Sprintf("like_video:%d", id)
}

// 读取用户点赞列表到Redis
func loadLikeListByUser(uid uint) error {
	// 从MySQL中获取数据
	likeList, err := dao.GetLikeListByUserID(uid)
	if err != nil {
		return err
	}
	uidKey := generateLikeKeyByUser(uid)

	// 遍历点赞列表并发添加到Redis中
	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()

	wg, err := sync.WaitGroup{}, nil
	for _, like := range likeList {
		wg.Add(1)
		go func(like model.Like) {
			defer wg.Done()
			if e := redis.RdbLike.SAdd(ctx, uidKey, like.VideoID).Err(); e != nil {
				err = e
			}
		}(*like)
	}
	wg.Wait()

	// 只要有一次写入发生错误，都要返回错误
	return err
}

// 读取视频点赞列表到Redis
func loadLikeListByVideo(vid uint) error {
	likeList, err := dao.GetLikeListByVideoID(vid)
	if err != nil {
		return err
	}
	vidKey := generateLikeKeyByVideo(vid)

	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()

	wg, err := sync.WaitGroup{}, nil
	for _, like := range likeList {
		wg.Add(1)
		go func(like model.Like) {
			defer wg.Done()
			if e := redis.RdbLike.SAdd(ctx, vidKey, like.UserID).Err(); e != nil {
				err = e
			}
		}(*like)
	}
	wg.Wait()

	return err
}
