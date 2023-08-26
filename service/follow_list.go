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

type FollowPair struct {
	Idx    int
	Follow *model.User
}

// GetFollowList 获取关注列表
func GetFollowList(id, curID uint) ([]*model.User, bool) {
	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()
	key := redis.GenerateFollowKey(id)

	// 1. 加载关注列表到redis
	if LoadFollowList(ctx, id) != nil {
		return nil, false
	}

	// 2. 从redis中获取关注列表
	follows, err := redis.RdbFollow.SMembers(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	// 3. 根据关注列表的用户ID获取用户信息
	followChan := make(chan FollowPair)
	status, nullIdx := true, len(follows)
	for i, follow := range follows {
		go func(i int, follow string) {
			// 跳过并记录空值的下标
			if follow == config.RedisValueOfNULL {
				nullIdx = i
				return
			}
			_fid, _ := strconv.ParseUint(follow, 10, 64)
			fid := uint(_fid)
			user, ok := GetUserByID(fid, curID)
			if !ok {
				status = false
			}
			followChan <- FollowPair{Idx: i, Follow: user}
		}(i, follow)
	}

	// 4. 通过管道获取结果，添加到结果切片中
	if len(follows) <= 1 {
		return nil, len(follows) == 1
	}
	// 因为缓存中有一个空值，所以结果的元素个数应-1
	followList := make([]*model.User, len(follows)-1)
	for i := 1; i < len(follows); i++ {
		follow := <-followChan
		idx := follow.Idx
		// 跳过空值
		if idx >= nullIdx {
			idx--
		}
		followList[idx] = follow.Follow
	}

	return followList, status
}

// IsFollow 判断 id 是否关注 followID
// 返回值第一个bool为是否关注，第二个bool为查询是否成功
func IsFollow(id, followID uint) (bool, bool) {
	ctx, cancel := redis.WithTimeoutContextBySecond(2)
	defer cancel()
	key := redis.GenerateFollowKey(id)

	// 由于SIsMember无法判断出key是否存在，只能判断成员是否存在
	// 而判断key是否存在需要执行Exists得到，所以至少需要执行两条redis命令
	// 为了使有缓存状态下的性能提升到极致，这里把这两条命令并发执行
	// 不过与此同时会增加无缓存状态下的redis命令数，所以无缓存状态下会慢一些
	wg := sync.WaitGroup{}
	wg.Add(2)
	var exists, isFollow bool

	// 判断是否存在这个key
	go func() {
		defer wg.Done()
		n, err := redis.RdbFollow.Exists(ctx, key).Result()
		exists = n == 1 && err == nil
	}()

	// 从redis获取数据返回
	go func() {
		defer wg.Done()
		f, err := redis.RdbFollow.SIsMember(ctx, key, followID).Result()
		if err != nil {
			exists = false
		}
		isFollow = f
	}()

	// 更新过期时间
	go func() {
		redis.RdbFollow.Expire(ctx, key, config.RedisKeyTTL)
	}()

	// 如果存在这个key，说明SIsMember查询出来的是有效数据，直接返回
	wg.Wait()
	if exists {
		return isFollow, true
	}

	// 否则从数据库加载
	if LoadFollowList(ctx, id) != nil {
		return false, false
	}

	// 重新获取数据
	isFollow, err := redis.RdbFollow.SIsMember(ctx, key, followID).Result()

	return isFollow, err == nil

}

// LoadFollowList 加载关注列表到redis
func LoadFollowList(ctx context.Context, id uint) error {
	key := redis.GenerateFollowKey(id)

	// 刷新过期时间
	defer func() {
		go func() { redis.RdbFollow.Expire(ctx, key, config.RedisKeyTTL) }()
	}()

	// 判断redis里是否存在关注列表，不存在时从MySQL读取到Redis
	if n, err := redis.RdbFollow.Exists(ctx, key).Result(); n == 0 || err != nil {
		// 从数据库读取数据
		followList, err := dao.GetBasicFollowList(id)
		if err != nil {
			return err
		}

		// 先添加一个默认值，防止缓存穿透
		if err := redis.RdbFollow.SAdd(ctx, key, config.RedisValueOfNULL).Err(); err != nil {
			return err
		}

		// 并发写入redis
		wg := sync.WaitGroup{}
		for _, follow := range followList {
			wg.Add(1)
			go func(follow model.Follow) {
				defer wg.Done()
				if e := redis.RdbFollow.SAdd(ctx, key, follow.UserID).Err(); e != nil {
					err = e
				}
			}(*follow)
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
