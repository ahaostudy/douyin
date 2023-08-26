package service

import (
	"context"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"strconv"
	"sync"
)

type FollowerPair struct {
	Idx      int
	Follower *model.User
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(id, curID uint) ([]*model.User, bool) {
	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()
	key := redis.GenerateFollowerKey(id)

	// 1. 加载粉丝列表到redis
	if LoadFollowerList(ctx, id) != nil {
		return nil, false
	}

	// 2. 从redis中获取粉丝列表
	followers, err := redis.RdbFollow.SMembers(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	// 3. 根据粉丝列表的用户ID获取用户信息
	followerChan := make(chan FollowerPair)
	status, nullIdx := true, len(followers)
	for i, follower := range followers {
		go func(i int, follower string) {
			// 跳过并记录空值的下标
			if follower == config.RedisValueOfNULL {
				nullIdx = i
				return
			}
			_fid, _ := strconv.ParseUint(follower, 10, 64)
			fid := uint(_fid)
			user, ok := GetUserByID(fid, curID)
			if !ok {
				status = false
			}
			followerChan <- FollowerPair{Idx: i, Follower: user}
		}(i, follower)
	}

	// 4. 通过管道获取结果，添加到结果切片中
	if len(followers) <= 1 {
		return nil, len(followers) == 1
	}

	followerList := make([]*model.User, len(followers)-1)
	for i := 1; i < len(followers); i++ {
		follower := <-followerChan
		idx := follower.Idx
		if idx >= nullIdx {
			idx--
		}
		followerList[idx] = follower.Follower
	}

	return followerList, status
}

// LoadFollowerList 加载粉丝列表到redis
func LoadFollowerList(ctx context.Context, id uint) error {
	key := redis.GenerateFollowerKey(id)

	// 刷新过期时间
	defer func() {
		go func() { redis.RdbFollow.Expire(ctx, key, config.RedisKeyTTL) }()
	}()

	// 判断redis里是否存在关注列表，不存在时从MySQL读取到Redis
	if n, err := redis.RdbFollow.Exists(ctx, key).Result(); n == 0 || err != nil {
		// 从数据库读取数据
		followerList, err := dao.GetBasicFollowerList(id)
		if err != nil {
			return err
		}

		// 先添加一个默认值，防止缓存穿透
		if err := redis.RdbFollow.SAdd(ctx, key, config.RedisValueOfNULL).Err(); err != nil {
			return err
		}

		// 并发写入redis
		wg := sync.WaitGroup{}
		for _, follower := range followerList {
			wg.Add(1)
			go func(follower model.Follow) {
				defer wg.Done()
				if e := redis.RdbFollow.SAdd(ctx, key, follower.FollowerID).Err(); e != nil {
					err = e
				}
			}(*follower)
		}
		wg.Wait()

		// 如果过程中读取失败，直接将key删除，防止脏写
		if err != nil {
			redis.RdbFollow.Del(ctx, key)
			return err
		}
	}

	return nil
}
